import {Footer} from '@/components';
import {login, register} from '@/services/short-link/user';
import {LockOutlined, UserOutlined,} from '@ant-design/icons';
import {LoginForm, ModalForm, ProFormText} from '@ant-design/pro-components';
import {FormattedMessage, Helmet, history, SelectLang, useIntl, useModel} from '@umijs/max';
import {Alert, message, Tabs} from 'antd';
import Settings from '../../../../config/defaultSettings';
import React, {useState} from 'react';
import {flushSync} from 'react-dom';
import {createStyles} from 'antd-style';


const useStyles = createStyles(({token}) => {
    return {
        action: {
            marginLeft: '8px',
            color: 'rgba(0, 0, 0, 0.2)',
            fontSize: '24px',
            verticalAlign: 'middle',
            cursor: 'pointer',
            transition: 'color 0.3s',
            '&:hover': {
                color: token.colorPrimaryActive,
            },
        },
        lang: {
            width: 42,
            height: 42,
            lineHeight: '42px',
            position: 'fixed',
            right: 16,
            borderRadius: token.borderRadius,
            ':hover': {
                backgroundColor: token.colorBgTextHover,
            },
        },
        container: {
            display: 'flex',
            flexDirection: 'column',
            height: '100vh',
            overflow: 'auto',
            backgroundImage:
                "url('https://mdn.alipayobjects.com/yuyan_qk0oxh/afts/img/V-_oS6r-i7wAAAAAAAAAAAAAFl94AQBr')",
            backgroundSize: '100% 100%',
        },
    };
});

const Lang = () => {
    const {styles} = useStyles();

    return (
        <div className={styles.lang} data-lang>
            {SelectLang && <SelectLang/>}
        </div>
    );
};

const LoginMessage: React.FC<{
    content: string;
}> = ({content}) => {
    return (
        <Alert
            style={{
                marginBottom: 24,
            }}
            message={content}
            type="error"
            showIcon
        />
    );
};

const Login: React.FC = () => {
    const [userLoginState, setUserLoginState] = useState<API.LoginResult>({});
    const [type, setType] = useState<string>('account');
    const {initialState, setInitialState} = useModel('@@initialState');
    const {styles} = useStyles();
    const intl = useIntl();

    const fetchUserInfo = async () => {
        const userInfo = await initialState?.fetchUserInfo?.();
        if (userInfo) {
            flushSync(() => {
                setInitialState((s) => ({
                    ...s,
                    currentUser: userInfo,
                }));
            });
        }
    };

    const handleSubmit = async (values: API.LoginParams) => {
        try {
            // 登录
            const msg = await login({...values, type});
            console.log(msg);
            if (msg.code === 200) {
                const defaultLoginSuccessMessage = intl.formatMessage({
                    id: 'pages.login.success',
                    defaultMessage: '登录成功！',
                });
                await fetchUserInfo();
                localStorage.setItem("token", msg.data.token)
                message.success(defaultLoginSuccessMessage);
                const urlParams = new URL(window.location.href).searchParams;
                history.push(urlParams.get('redirect') || '/');
                return;
            }
            message.error(msg.msg)
            setUserLoginState(msg);
        } catch (error) {
            const defaultLoginFailureMessage = intl.formatMessage({
                id: 'pages.login.failure',
                defaultMessage: '登录失败，请重试！',
            });
            console.log(error);
            message.error(defaultLoginFailureMessage);
        }
    };
    const {code, type: loginType} = userLoginState;

    return (
        <div className={styles.container}>
            <Helmet>
                <title>
                    {intl.formatMessage({
                        id: 'menu.login',
                        defaultMessage: '登录页',
                    })}
                    - {Settings.title}
                </title>
            </Helmet>
            <Lang/>
            <div
                style={{
                    flex: '1',
                    padding: '32px 0',
                }}
            >
                <LoginForm
                    contentStyle={{
                        minWidth: 280,
                        maxWidth: '75vw',
                    }}
                    logo={<img alt="logo" src="/logo.svg"/>}
                    title="Ant Design"
                    subTitle={intl.formatMessage({id: 'pages.layouts.userLayout.title'})}
                    initialValues={{
                        autoLogin: true,
                    }}

                    onFinish={async (values) => {
                        await handleSubmit(values as API.LoginParams);
                    }}
                >
                    <Tabs
                        activeKey={type}
                        onChange={setType}
                        centered
                        items={[
                            {
                                key: 'account',
                                label: intl.formatMessage({
                                    id: 'pages.login.accountLogin.tab',
                                    defaultMessage: '账户密码登录',
                                }),
                            },

                        ]}
                    />

                    {code !== 200 && loginType === 'account' && (
                        <LoginMessage
                            content={intl.formatMessage({
                                id: 'pages.login.accountLogin.errorMessage',
                                defaultMessage: '账户或密码错误(admin/ant.design)',
                            })}
                        />
                    )}
                    {type === 'account' && (
                        <>
                            <ProFormText
                                value={"admin-2"}
                                name="username"
                                fieldProps={{
                                    size: 'large',
                                    prefix: <UserOutlined/>,
                                }}
                                placeholder={intl.formatMessage({
                                    id: 'pages.login.username.placeholder',
                                    defaultMessage: '用户名',
                                })}
                                rules={[
                                    {
                                        required: true,
                                        message: (
                                            <FormattedMessage
                                                id="pages.login.username.required"
                                                defaultMessage="请输入用户名!"
                                            />
                                        ),
                                    },
                                ]}
                            />
                            <ProFormText.Password
                                value={"Wanyuqin1@"}
                                name="password"
                                fieldProps={{
                                    size: 'large',
                                    prefix: <LockOutlined/>,
                                }}
                                placeholder={intl.formatMessage({
                                    id: 'pages.login.password.placeholder',
                                    defaultMessage: '密码',
                                })}
                                rules={[
                                    {
                                        required: true,
                                        message: (
                                            <FormattedMessage
                                                id="pages.login.password.required"
                                                defaultMessage="请输入密码！"
                                            />
                                        ),
                                    },
                                ]}
                            />
                        </>
                    )}

                    <div
                        style={{
                            marginBottom: 24,
                        }}
                    >

                        <ModalForm
                            onFinish={async (values) => {
                                const ret = await register(values)
                                if (ret.code == 200) {
                                    message.success("注册成功")
                                } else {
                                    message.error(ret.msg)
                                    return false
                                }
                                return true;
                            }}
                            modalProps={{
                                destroyOnClose: true
                            }}
                            title="注册用户"
                            trigger={<a style={{float: 'right',}}>
                                用户注册
                            </a>}

                        >
                            <ProFormText
                                width={"md"}
                                name={"username"}
                                label={"用户名"}
                                tooltip={"用户名长度不能小于6且不能大于30"}
                                rules={[{required: true}]}
                                placeholder={"请输入用户名"}>
                            </ProFormText>
                            <ProFormText
                                width={"md"}
                                name={"password"}
                                label={"密码"}
                                tooltip={"密码最少8位包含大小写、数字、特殊字符"}
                                rules={[{required: true}]}
                                placeholder={"请输入密码"}>
                            </ProFormText>
                        </ModalForm>
                    </div>
                </LoginForm>
            </div>
            <Footer/>
        </div>
    );
};

export default Login;
