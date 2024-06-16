import React from 'react';
import {ModalForm, ProFormDateTimePicker, ProFormText, ProTable, TableDropdown} from '@ant-design/pro-components';
import {Button,message} from 'antd';
import {PlusOutlined} from '@ant-design/icons';
import {addLink} from '@/services/short-link/link';
import moment from 'moment';


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
                查看
            </a>,
            <TableDropdown
                key="actionGroup"
                onSelect={() => action?.reload()}
                menus={[
                    {key: 'copy', name: '复制'},
                    {key: 'delete', name: '删除'},
                ]}
            />,
        ],
    }
]


const LinkList: React.FC = () => {
    return <ProTable
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
