'use strict'

/********************************************************************
 created:    2023-03-22
 author:     lixianmin

 Copyright (C) - All Rights Reserved
 *********************************************************************/

// 将[Uint8Array(), Uint8Array()] 合并为一个单一的Uint8Array()
export function joinUint8Array(list) {
    if (Array.isArray(list)) {
        let length = 0
        list.forEach(item => {
            length += item.length
        })

        const data = new Uint8Array(length)
        let offset = 0
        list.forEach(item => {
            data.set(item, offset)
            offset += item.length
        })

        return data
    }
}