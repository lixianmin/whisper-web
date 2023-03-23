'use strict'
import {onMount} from "solid-js";

/********************************************************************
 created:    2023-03-20
 author:     lixianmin

 Copyright (C) - All Rights Reserved
 *********************************************************************/

export default function () {
    // todo UI控制相关的，应该跟model层拆到两个地方
    onMount(() => {
        const testKey = 'Control'
        document.addEventListener('keydown', evt => {
            if (evt.key === testKey) {
                console.log('start recording~')
                startRecording()
            }
        })

        document.addEventListener('keyup', evt => {
            if (evt.key === testKey) {
                stopRecording()
                console.log('stop recording~')
            }
        })
    })

    const kMaxAudioSeconds = 120
    const kSampleRate = 16000

    const audioContext = new AudioContext({
        sampleRate: kSampleRate,
        channelCount: 1,
        echoCancellation: false,
        autoGainControl: true,
        noiseSuppression: true,
    })

    let isRecording = false
    let mediaRecorder = undefined

    // record up to kMaxAudio_s seconds of audio from the microphone
    // check if doRecording is false every 1000 ms and stop recording if so
    // update progress information
    function startRecording() {
        if (isRecording) {
            return
        }

        isRecording = true

        navigator.mediaDevices.getUserMedia({audio: true, video: false})
            .then(function (stream) {
                mediaRecorder = new MediaRecorder(stream)
                let chunks = []
                mediaRecorder.ondataavailable = function (evt) {
                    chunks.push(evt.data)
                }

                function stopTracks() {
                    stream.getTracks().forEach(function (track) {
                        track.stop()
                    })
                }

                // 录音结束的时候，会调用onstop，然后把chunks中的内容写到blob中，而后是使用reader读取blob中的内容，读成功后走到reader.onload()
                mediaRecorder.onstop = function (e) {
                    const fileReader = new FileReader()
                    fileReader.onload = function () {
                        const result = fileReader.result
                        transcribeAudio(result)
                    }

                    const blob = new Blob(chunks, {'type': 'audio/ogg; codecs=opus'})
                    chunks = []
                    fileReader.readAsArrayBuffer(blob)
                    stopTracks()
                }

                mediaRecorder.start()
            })
            .catch(function (err) {
                console.error('error getting audio stream: ' + err)
            })
    }

    function stopRecording() {
        if (isRecording) {
            isRecording = false
            mediaRecorder.stop()
        }
    }

    function transcribeAudio(data) {
        // const blob = new Blob(chunks, {'type': 'audio/ogg; codecs=opus'})
        const buf = new Uint8Array(data)
        audioContext.decodeAudioData(buf.buffer, function (audioBuffer) {
            const offlineContext = new OfflineAudioContext(audioBuffer.numberOfChannels, audioBuffer.length, audioBuffer.sampleRate)
            const source = offlineContext.createBufferSource()
            source.buffer = audioBuffer
            source.connect(offlineContext.destination)
            source.start(0)

            offlineContext.startRendering().then(function (renderedBuffer) {
                let audio = renderedBuffer.getChannelData(0)
                console.log('audio recorded, size: ' + audio.length)

                // truncate to first 30 seconds
                if (audio.length > kMaxAudioSeconds * kSampleRate) {
                    audio = audio.slice(0, kMaxAudioSeconds * kSampleRate);
                    console.log('truncated audio to first ' + kMaxAudioSeconds + ' seconds')
                }
                onSetAudio(audio)
            })
        }, function (e) {
            console.error('error decoding audio: ' + e)
        })
    }

    // 录音结束了，设置到这里
    let instance = undefined

    function onSetAudio(audio) {
        if (!instance) {
            instance = Module.init('whisper.bin');
        }

        // printText('')
        // printText('js: processing - this might take a while ...')
        // printText('')

        setTimeout(function () {
            const language = 'zh'
            const nthreads = 8
            const translate = false

            const ret = Module.full_default(instance, audio, language, nthreads, translate)
            // console.log('js: full_default returned: ' + ret)
            if (ret) {
                // printText("js: whisper returned: " + ret)
            }
        }, 100)
    }

    return <>
    </>
}