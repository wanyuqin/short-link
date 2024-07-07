import React, {useEffect, useRef, useState} from 'react';
import {Button, Drawer, InputNumber, message, Modal, Space} from 'antd';
import {ModalForm, ProTable} from '@ant-design/pro-components';
import {addBlackList, delBlackList, listBlackList, updateBlackList} from "@/services/short-link/blacklist";
import {PlusOutlined} from "@ant-design/icons";

// 定义传入值
interface BlacklistDrawerProps {
    visible: boolean;
    onClose: () => void;
    shortUrl: string;
}

const onChange: InputNumberProps['onChange'] = (value) => {
    console.log('changed', value);
};

const columns = [
    {
        title: "id",
        dataIndex: "id",
        search: false
    },
    {
        title: "IP",
        dataIndex: "ip",
    },
    {
        title: "状态",
        dataIndex: "status",
        onFilter: true,
        ellipsis: true,
        valueEnum: {
            1: {text: '启用', status: 'Success'},
            0: {text: '禁用', status: 'Error'},
        },
    },
    {
        title: "创建时间",
        dataIndex: "createdAt",
        search: false
    },
    {
        title: '操作',
        valueType: 'option',
        key: 'option',
        render: (text, record, _, action) => [
            <a onClick={async () => {
                Modal.confirm({
                    title: '确认删除',
                    content: '确定要删除这个链接吗？',
                    okText: '确认',
                    cancelText: '取消',
                    onOk: async () => {
                        console.log(record?.id)
                        const res = await delBlackList({"id": record?.id});
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
            <a onClick={async () => {
                // 根据状态执行启用或禁用操作
                const newStatus = record.status === 0 ? 1 : 0;
                const res = await updateBlackList({"id": record?.id, "status": newStatus, "shortUrl": record.shortUrl});
                if (res.code === 200) {
                    message.success(newStatus === 0 ? "已启用" : "已禁用");
                    action?.reloadAndRest();
                } else {
                    message.error(res.msg);
                }
            }}>
                {record.status === 0 ? "启用" : "禁用"}
            </a>
        ],
    }
]

const BlacklistDrawer: React.FC<BlacklistDrawerProps> = ({visible, onClose, shortUrl}) => {
    const actionRef = useRef<any>();
    const [IPPart1, setIPPart1] = useState(0);
    const [IPPart2, setIPPart2] = useState(0);
    const [IPPart3, setIPPart3] = useState(0);
    const [IPPart4, setIPPart4] = useState(0);
    // 重置 IP 地址部分的状态
    const resetIPParts = () => {
        setIPPart1(0);
        setIPPart2(0);
        setIPPart3(0);
        setIPPart4(0);
    };

    // 黑名单列表
    const fetchBlackList = async (params) => {
        const {pageSize, page, IP} = params
        const res = await listBlackList({page, pageSize, IP, shortUrl})
        const {data: nestedData, total: total} = res.data;
        return {
            data: nestedData,
            total: total,
            success: res.code === 200
        }

    }


    useEffect(() => {
        if (visible && actionRef.current) {
            console.log(shortUrl)
            actionRef.current.reload();
        }
    }, [visible, shortUrl]);

    return (
        <Drawer
            title="黑名单"
            width={720}
            onClose={onClose}
            visible={visible}
            footer={
                <div
                    style={{
                        textAlign: 'right',
                    }}
                >
                    <Button onClick={onClose} style={{marginRight: 8}}>
                        关闭
                    </Button>
                </div>
            }
        >
            <ProTable
                actionRef={actionRef}
                rowKey="id"
                columns={columns}
                pagination={{
                    pageSize: 10,
                }}
                headerTitle="黑名单详情"
                request={async (params: T & {
                    pageSize: number;
                    current: number;
                }) => {
                    const response = await fetchBlackList({
                        page: params.current,
                        pageSize: params.pageSize,
                    })
                    console.log(response)
                    return response
                }}
                toolBarRender={() => [
                    <ModalForm
                        trigger={
                            <Button type="primary">
                                <PlusOutlined/>
                                新建
                            </Button>
                        }
                        onFinish={async () => {
                            const IP = IPPart1 + "." + IPPart2 + "." + IPPart3 + "." + IPPart4
                            const ret = await addBlackList({IP, shortUrl})
                            console.log(ret)
                            if (ret.code == 200) {
                                message.success("添加成功");
                                actionRef.current?.reload(); // 刷新表格数据
                                resetIPParts(); // 重置 IP 地址部分的状态
                                return true;
                            }
                            message.error(ret.msg);
                            return false;

                        }}
                        modalProps={{
                            onCancel: resetIPParts // 当取消时也重置 IP 地址部分的状态
                        }}
                    >
                        <Space wrap>
                            IP地址
                            <InputNumber size="small" min={0} value={IPPart1} max={255} onChange={(value) => {
                                setIPPart1(value)
                            }}/>.
                            <InputNumber size="small" min={0} value={IPPart2} max={255} onChange={(value) => {
                                setIPPart2(value)
                            }}/>.
                            <InputNumber size="small" min={0} value={IPPart3} max={255} onChange={(value) => {
                                setIPPart3(value)
                            }}/>.
                            <InputNumber size="small" min={0} value={IPPart4} max={255} onChange={(value) => {
                                setIPPart4(value)
                            }}/>
                        </Space>
                    </ModalForm>
                ]}

            />
        </Drawer>


    );
};

export default BlacklistDrawer;
