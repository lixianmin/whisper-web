'use strict'
import Dexie from "dexie";

/********************************************************************
 created:    2023-03-22
 author:     lixianmin

 Copyright (C) - All Rights Reserved
 *********************************************************************/

export function createWisperModel() {
    const db = new Dexie("whisper_model")
    db.version(1).stores({
        models: '++id, name, chunk',
    })

    async function _isModelExists(name) {
        if (name && db.models) {
            const first = await db.models.where({name}).first()
            if (first) {
                return true
            }
        }

        return false
    }

    async function _saveModel(name, data) {
        if (!name || !data) {
            return
        }

        const chunkSize = 1024 * 1024 * 100
        const count = data.length / chunkSize

        db.transaction('rw', db.models, async () => {
            const exists = await _isModelExists(name)
            if (!exists) {
                for (let i = 0; i < count; i++) {
                    const chunk = data.subarray(i * chunkSize, (i + 1) * chunkSize)
                    await db.models.add({name, chunk})
                }
            }
        })
    }

    async function _loadModel(name) {
        if (db.models) {
            const chunks = await db.models.where({name}).toArray()
            const data = chunks.join()
            return data
        }
    }

    return {
        isModelExists: _isModelExists,
        saveModel: _saveModel,
        loadModel: _loadModel,
    }
}