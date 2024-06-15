// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 添加链接 添加链接 POST /link/add */
export async function postLinkAdd(body: API.AddLinkReq, options?: { [key: string]: any }) {
  return request<API.Response>('/link/add', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
