import {PageContainer} from '@ant-design/pro-components';
import {useModel} from '@umijs/max';
import {Card, theme, Typography} from 'antd';
import React from 'react';

const {Title} = Typography;

/**
 * 每个单独的卡片，为了复用样式抽成了组件
 * @param param0
 * @returns
 */
const InfoCard: React.FC<{
    title: string;
    index: number;
    desc: string;
}> = ({title, index, desc}) => {
    const {useToken} = theme;

    const {token} = useToken();

    return (
        <div
            style={{
                backgroundColor: token.colorBgContainer,
                boxShadow: token.boxShadow,
                borderRadius: '8px',
                fontSize: '14px',
                color: token.colorTextSecondary,
                lineHeight: '22px',
                padding: '16px 19px',
                minWidth: '220px',
                flex: 1,
            }}
        >
            <div
                style={{
                    display: 'flex',
                    gap: '4px',
                    alignItems: 'center',
                }}
            >
                <div
                    style={{
                        width: 48,
                        height: 48,
                        lineHeight: '22px',
                        backgroundSize: '100%',
                        textAlign: 'center',
                        padding: '8px 16px 16px 12px',
                        color: '#FFF',
                        fontWeight: 'bold',
                        backgroundImage:
                            "url('https://gw.alipayobjects.com/zos/bmw-prod/daaf8d50-8e6d-4251-905d-676a24ddfa12.svg')",
                    }}
                >
                    {index}
                </div>
                <div
                    style={{
                        fontSize: '16px',
                        color: token.colorText,
                        paddingBottom: 8,
                    }}
                >
                    {title}
                </div>
            </div>
            <div
                style={{
                    fontSize: '14px',
                    color: token.colorTextSecondary,
                    textAlign: 'justify',
                    lineHeight: '22px',
                    marginBottom: 8,
                }}
            >
                {desc}
            </div>
        </div>
    );
};

const Welcome: React.FC = () => {
    const {token} = theme.useToken();
    const {initialState} = useModel('@@initialState');
    return (
        <PageContainer>
            <Card
                style={{
                    borderRadius: 8,
                }}
                bodyStyle={{
                    backgroundImage:
                        initialState?.settings?.navTheme === 'realDark'
                            ? 'background-image: linear-gradient(75deg, #1A1B1F 0%, #191C1F 100%)'
                            : 'background-image: linear-gradient(75deg, #FBFDFF 0%, #F5F7FF 100%)',
                }}
            >
                <div
                    style={{
                        backgroundPosition: '100% -30%',
                        backgroundRepeat: 'no-repeat',
                        backgroundSize: '274px auto',
                        backgroundImage:
                            "url('https://gw.alipayobjects.com/mdn/rms_a9745b/afts/img/A*BuFmQqsB2iAAAAAAAAAAAAAAARQnAQ')",
                    }}
                >
                    <div
                        style={{
                            fontSize: '20px',
                            color: token.colorTextHeading,
                        }}
                    >
                        欢迎了解-使用short-link
                    </div>
                    <p
                        style={{
                            fontSize: '14px',
                            color: token.colorTextSecondary,
                            lineHeight: '22px',
                            marginTop: 16,
                            marginBottom: 32,
                            width: '65%',
                        }}
                    >
                        <Typography.Title level={2} style={{margin: 0}}>
                            概述
                        </Typography.Title>
                        Short-Link 是一个功能强大且易于使用的短链接生成和管理平台。它旨在帮助用户创建、分享和跟踪短链接，以提高
                        URL 的可读性和用户体验。无论是个人用户还是企业用户，Short-Link 都提供了丰富的功能来满足不同需求。
                        <Typography.Title level={2} style={{margin: 0}}>
                            功能特性
                        </Typography.Title>

                    </p>
                    <div
                        style={{
                            display: 'flex',
                            flexWrap: 'wrap',
                            gap: 16,
                        }}
                    >
                        <InfoCard
                            index={1}
                            title="短链接生成"
                            desc="Short-Link 支持快速生成简洁的短链接。用户只需输入原始 URL，平台会自动生成唯一的短链接，便于分享和传播。"
                        />
                        <InfoCard
                            index={2}
                            title="访问统计"
                            desc="通过详尽的统计功能，用户可以实时监控短链接的访问情况。包括访问次数、访客来源、访问时间、地理位置等，帮助用户深入了解链接的使用情况。"
                        />
                        <InfoCard
                            index={3}
                            title="自定义短链接"
                            desc="用户可以根据需要自定义短链接的后缀，使其更符合个人或品牌的需求。"
                        />


                    </div>

                    <div
                        style={{
                            display: 'flex',
                            flexWrap: 'wrap',
                            gap: 16,
                        }}
                    >
                        <InfoCard
                            index={4}
                            title="链接管理"
                            desc="提供便捷的管理界面，用户可以轻松地编辑、删除和查看已创建的短链接，并批量处理多个链接。"
                        />
                        <InfoCard
                            index={5}
                            title="安全与隐私"
                            desc="为了确保短链接的安全性，Short-Link 采用了多种安全措施，包括链接有效期设置、访问权限控制和防止恶意点击等。"
                        />
                        <InfoCard
                            index={5}
                            title="GITHUB"
                            desc="为了确保短链接的安全性，Short-Link 采用了多种安全措施，包括链接有效期设置、访问权限控制和防止恶意点击等。"
                        />
                    </div>
                </div>
            </Card>
        </PageContainer>
    );
};

export default Welcome;
