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


    function storeFS(fname, buf) {
        // write to WASM file using FS_createDataFile
        // if the file exists, delete it
        try {
            Module.FS_unlink(fname);
        } catch (e) {
            // ignore
        }

        Module.FS_createDataFile("/", fname, buf, true, true);

        printTextarea('storeFS: stored model: ' + fname + ' size: ' + buf.length);
    }

    function loadWhisper() {
        let url     = 'https://localhost/ggml-model-whisper-base.bin'
        let dst     = 'whisper.bin';
        let size_mb = 142;

        let cbProgress = function(p) {
            // console.log(p)
        };

        const cbCancel = function() {
            console.log('canceled')
        }

        loadRemote(url, dst, size_mb, cbProgress, storeFS, cbCancel, printTextarea);
    }

    loadWhisper()

    return <>
        <AudioRecorder/>
    </>
}

const root = document.getElementById('root')
render(() => <App />, root)