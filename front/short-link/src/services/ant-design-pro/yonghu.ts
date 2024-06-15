// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 用户登陆 用户登陆 POST /users/login */
export async function postUsersLogin(body: API.LoginReq, options?: { [key: string]: any }) {
  return request<API.Response>('/users/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 用户注册 用户注册 POST /users/register */
export async function postUsersRegister(body: API.RegisterReq, options?: { [key: string]: any }) {
  return request<API.Response>('/users/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
