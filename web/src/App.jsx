/********************************************************************
 created:    2023-03-21
 author:     lixianmin

 Copyright (C) - All Rights Reserved
 *********************************************************************/
import './App.module.css';
import AudioRecorder from "./widgets/AudioRecorder";
import {render} from "solid-js/web";

function App() {
    return <>
        <AudioRecorder/>
    </>
}

const root = document.getElementById('root')
render(() => <App />, root)