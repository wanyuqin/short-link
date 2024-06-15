declare namespace API {
  type AddLinkReq = {
    expiredAt?: number;
    originalUrl?: string;
    userId?: number;
  };

  type LoginReq = {
    password?: string;
    username?: string;
  };

  type RegisterReq = {
    /** 密码 */
    password?: string;
    /** 用户名 */
    username?: string;
  };

  type Response = {
    code?: number;
    data?: any;
    msg?: string;
    timeStamp?: number;
  };
}
