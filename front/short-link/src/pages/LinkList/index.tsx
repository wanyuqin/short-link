import React, {useRef} from 'react';
import {ModalForm, ProFormDateTimePicker, ProFormText, ProTable} from '@ant-design/pro-components';
import {Button, message} from 'antd';
import {PlusOutlined} from '@ant-design/icons';
import {addLink, delLink, linkList} from '@/services/short-link/link';

const fetchLinkList = async () => {
    const res = await linkList({})
    const {data: nestedData, lastId} = res.data;
    return {
        data: nestedData,
        success: res.code === 200,
        total: lastId,
    };
}

const columns = [
    {
        title: "原链接",
        dataIndex: "originUrl",

    },
    {
        title: "短链接",
        dataIndex: "shortUrl",

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
                onClick={() => {
                    action?.startEditable?.(record.id);
                }}
            >
                编辑
            </a>,
            <a href={record.url} target="_blank" rel="noopener noreferrer" key="view">
                复制
            </a>,
            <a onClick={async () => {
                const res = await delLink({
                    "shortUrl": record?.shortUrl
                })
                if (res.code === 200) {
                    message.success("删除成功");
                    action?.reload(); // 刷新表格数据
                } else {
                    message.error(res.msg);
                }
            }}>
                删除
            </a>,
            // <TableDropdown
            //     key="actionGroup"
            //     onSelect={() => action?.reload()}
            //     menus={[
            //         {key: 'copy', name: '复制'},
            //         {key: 'delete', name: '删除'},
            //     ]}
            // />,
        ],
    }
]


const LinkList: React.FC = () => {
    const actionRef = useRef<any>(); // 创建一个ref来存储ProTable的实例

    return <ProTable
        actionRef={actionRef}
        request={fetchLinkList}
        columns={columns}
        toolBarRender={() => [
            <ModalForm
                trigger={
                    <Button type="primary">
                        <PlusOutlined/>
                        新建
                    </Button>
                }
                onFinish={async (values) => {
                    const ret = await addLink(values)
                    if (ret.code == 200) {
                        message.success("添加成功")
                        actionRef.current?.reload(); // 刷新表格数据
                        return true
                    }
                    message.error(ret.msg)
                    return false
                }}
                modalProps={{
                    destroyOnClose: true
                }}
                title={"新建短链"}

            >
                <ProFormText
                    width="md"
                    name="originUrl"
                    label="原始链接"
                    tooltip=""
                    placeholder="请输入原始链接"
                />

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
        ]}>

    </ProTable>
}

export default LinkList;
