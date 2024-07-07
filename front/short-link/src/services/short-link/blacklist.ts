import {request} from "@umijs/max";

export async function addBlackList(data) {
  return request("http://127.0.0.1:8088/short-link/api/v1/admin/link/black-list/add", {
    method: 'POST',
    data
  })
}


export async function listBlackList(data) {
  return request("http://127.0.0.1:8088/short-link/api/v1/admin/link/black-list/list", {
    method: 'POST',
    data
  })
}

export async function delBlackList(data) {
  return request("http://127.0.0.1:8088/short-link/api/v1/admin/link/black-list/del", {
    method: 'POST',
    data
  })
}

export async function updateBlackList(data) {
  return request("http://127.0.0.1:8088/short-link/api/v1/admin/link/black-list/update-status", {
    method: 'POST',
    data
  })
}
