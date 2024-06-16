import {request} from "@umijs/max";


export async function addLink(data) {
    return request("http://127.0.0.1:8088/api/v1/admin/link/add", {
        method: 'POST',
        data
    })
}
