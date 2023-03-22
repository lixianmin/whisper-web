// Common Javascript functions used by the examples

import {createWisperModel} from "./dao/whisper_model"
import {joinUint8Array} from "./core/tools";

// fetch a remote file from remote URL using the Fetch API
async function fetchRemote(url, cbPrint) {
    cbPrint('fetchRemote: downloading with fetch()...');

    const response = await fetch(
        url,
        {
            method: 'GET',
            headers: {
                'Content-Type': 'application/octet-stream',
            },
        }
    );

    if (!response.ok) {
        cbPrint('fetchRemote: failed to fetch ' + url)
        return
    }

    const contentLength = response.headers.get('content-length')
    const total = parseInt(contentLength, 10)
    const reader = response.body.getReader()

    const chunks = []
    let receivedLength = 0
    let progressLast = -1

    while (true) {
        const {done, value} = await reader.read()

        if (done) {
            break
        }

        chunks.push(value)
        receivedLength += value.length

        if (contentLength) {
            const progressCur = Math.round((receivedLength / total) * 10)
            if (progressCur !== progressLast) {
                cbPrint('fetchRemote: fetching ' + 10 * progressCur + '% ...')
                progressLast = progressCur
            }
        }
    }

    const chunksAll = joinUint8Array(chunks)
    return chunksAll
}

// load remote data
// - check if the data is already in the IndexedDB
// - if not, fetch it from the remote URL and store it in the IndexedDB
export function loadRemote(url, dstFileName, cbReady, cbPrint) {
    navigator.storage.estimate().then(function (estimate) {
        cbPrint('loadRemote: storage quota: ' + estimate.quota + ' bytes')
        cbPrint('loadRemote: storage usage: ' + estimate.usage + ' bytes')
    });

    const whisperModel = createWisperModel()
    whisperModel.isModelExists(url).then(exists => {
        if (exists) {
            whisperModel.loadModel(url).then(data => {
                cbReady(dstFileName, data)
                cbPrint(`load from db: ${data.length}`)
            })
        } else {
            fetchRemote(url, cbPrint).then(data => {
                if (data) {
                    whisperModel.saveModel(url, data).then(() => {
                        cbReady(dstFileName, data)
                        console.log(`download from web: ${data.length}`)
                    })
                }
            })
        }
    })
}
