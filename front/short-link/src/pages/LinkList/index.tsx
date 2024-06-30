import React, {useRef, useState} from 'react';
import {
    ModalForm,
    ProFormDateTimePicker,
    ProFormFieldSet,
    ProFormSelect,
    ProFormText,
    ProTable
} from '@ant-design/pro-components';
import {Button, message, Modal} from 'antd';
import {PlusOutlined} from '@ant-design/icons';
import BlacklistDrawer from "@/pages/LinkList/components/BlacklistDrawer";
import {addLink, delLink, linkList} from '@/services/short-link/link';


const fetchLinkList = async (params) => {
    const {pageSize, page = 1, originUrl, lastId} = params;
    const res = await linkList({page, pageSize, originUrl});
    const {data: nestedData, total: total} = res.data;
    return {
        data: nestedData,
        success: res.code === 200,
        total: total,
    };
};

const LinkList: React.FC = () => {

    const columns = [
        {
            title: "原链接",
            dataIndex: "originUrl",
            ellipsis: true,
            tooltip: '原链接过长会自动收缩',
        },
        {
            title: "短链接",
            dataIndex: "shortUrl",
            search: false,
        },
        {
            title: "有效期",
            dataIndex: "expiredAt",
            search: false,
        },
        {
            title: "创建时间",
            dataIndex: "createdAt",
            search: false,
        },
        {
            title: '操作',
            valueType: 'option',
            key: 'option',
            render: (text, record, _, action) => [
                <a
                    key="editable"
                    onClick={() => handleBlacklistClick(record.shortUrl)}
                >
                    黑名单
                </a>,
                <a href={record.url} target="_blank" rel="noopener noreferrer" key="view">
                    复制
                </a>,
                <a onClick={async () => {
                    Modal.confirm({
                        title: '确认删除',
                        content: '确定要删除这个链接吗？',
                        okText: '确认',
                        cancelText: '取消',
                        onOk: async () => {
                            const res = await delLink({"shortUrl": record?.shortUrl});
                            if (res.code === 200) {
                                message.success("删除成功");
                                // 在这里重新加载表格数据，使用 initialLastId
                                action?.reloadAndRest();
                            } else {
                                message.error(res.msg);
                            }
                        }
                    });
                }}>
                    删除
                </a>,
            ],
        }
    ];

    const actionRef = useRef<any>(); // 创建一个ref来存储ProTable的实例
    const [lastId, setLastId] = useState(0); // 初始化 lastId 为 0
    const [initialLastId] = useState(0); // 保存初始的 lastId
    const [drawerVisible, setDrawerVisible] = useState(false);
    const [currentShortUrl, setCurrentShortUrl] = useState<string | null>(null);
    const handleBlacklistClick = (shortUrl: string) => {
        setCurrentShortUrl(shortUrl);
        setDrawerVisible(true);
    };
    return (
        <>
            <ProTable
                actionRef={actionRef}
                request={async (params, sorter, filter) => {
                    const response = await fetchLinkList({...params, page: params.current || 1, lastId});
                    console.log(response);
                    return response;
                }}
                pagination={{
                    pageSize: 5,
                }}
                columns={columns}
                rowKey="shortUrl"
                toolBarRender={() => [
                    <ModalForm
                        trigger={
                            <Button type="primary">
                                <PlusOutlined/>
                                新建
                            </Button>
                        }
                        onFinish={async (values) => {
                            console.log(values)
                            const ret = await addLink(values);
                            if (ret.code == 200) {
                                message.success("添加成功");
                                actionRef.current?.reload(); // 刷新表格数据
                                return true;
                            }
                            message.error(ret.msg);
                            return false;
                        }}
                        modalProps={{
                            destroyOnClose: true
                        }}
                        title="新建短链"
                    >
                        <ProFormFieldSet
                            name="originUrl"
                            label="原始链接"
                            type="group"
                            rules={[
                                {
                                    validator: (_, value) => {
                                        const [scheme, path] = value || [];
                                        if (!scheme) {
                                            return Promise.reject(new Error('协议不能为空'));
                                        }
                                        if (!path) {
                                            return Promise.reject(new Error('地址不能为空'));
                                        }
                                        return Promise.resolve();
                                    },
                                },
                            ]}
                            transform={(value: any) => {
                                const [scheme, path] = value || [];
                                return {originUrl: scheme + path};
                            }}
                        >
                            <ProFormSelect
                                name="scheme"
                                label="scheme"
                                valueEnum={{
                                    "https://": "https://",
                                    "http://": "http://",
                                }}></ProFormSelect>
                            <ProFormText width="md"/>
                        </ProFormFieldSet>
                        <ProFormDateTimePicker
                            width="md"
                            name="expiredAt"
                            label="有效期"
                            placeholder="请输入有效期"
                            fieldProps={{
                                format: (value) => value.format('YYYY-MM-DD hh:mm:ss'),
                            }}
                        />
                    </ModalForm>,
                ]}
                onReset={() => setLastId(initialLastId)} // 重置 lastId 为初始值
            ></ProTable>,

            <BlacklistDrawer
                visible={drawerVisible}
                onClose={() => setDrawerVisible(false)}
                shortUrl={currentShortUrl}
            />
        </>

    );
}

export default LinkList;
