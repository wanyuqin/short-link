import {request} from "@umijs/max";


export async function addLink(data) {
    return request("http://127.0.0.1:8088/short-link/api/v1/admin/link/add", {
        method: 'POST',
        data
    })
}


export async function linkList(data) {
    return request("http://127.0.0.1:8088/short-link/api/v1/admin/link/list", {
        method: 'POST',
        data
    })
}

export async function delLink(data){
    return request("http://127.0.0.1:8088/short-link/api/v1/admin/link/del", {
        method: 'POST',
        data
    })
}
