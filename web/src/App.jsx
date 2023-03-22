/********************************************************************
 created:    2023-03-21
 author:     lixianmin

 Copyright (C) - All Rights Reserved
 *********************************************************************/
import './App.module.css';
import AudioRecorder from "./widgets/AudioRecorder";
import {render} from "solid-js/web";
import {loadRemote} from "./code/helpers";

function App() {
    function storeFS(filename, buf) {
        // write to WASM file using FS_createDataFile
        // if the file exists, delete it
        try {
            Module.FS_unlink(filename)
        } catch (e) {
            // ignore
        }

        Module.FS_createDataFile("/", filename, buf, true, true)
        printText('storeFS: stored model: ' + filename + ' size: ' + buf.length)
    }

    function loadWhisper() {
        let url = 'https://localhost/ggml-model-whisper-base.bin'
        let dstFileName = 'whisper.bin'
        loadRemote(url, dstFileName, storeFS, printText)
    }

    loadWhisper()

    return <>
        <AudioRecorder/>
    </>
}

const root = document.getElementById('root')
render(() => <App/>, root)