import React, {useState} from 'react';
import type {MenuProps} from 'antd';
import {Button, DatePicker, Dropdown, Flex, message} from 'antd';
import {DownOutlined, SyncOutlined} from "@ant-design/icons";
import moment from 'moment';


const {RangePicker} = DatePicker;


const TimeSelector: React.FC<{ onTimeRangeChange: (dates: [moment.Moment, moment.Moment] | null) => void }> = ({onTimeRangeChange}) => {

    const [selectedOption, setSelectedOption] = useState<string | null>(null);
    const onClick: MenuProps['onClick'] = ({key, label}) => {
        const selectedItem = items.find((item) => item.key === key);
        if (selectedItem) {
            setSelectedOption(selectedItem.label as string);
            onTimeRangeChange()
        }
        // 设置

    };
    const items: MenuProps['items'] = [
        {key: '1', label: '近5分钟'},
        {key: '2', label: '近15分钟'},
        {key: '3', label: '近30分钟'},
        {key: '4', label: '近1小时'},
        {key: '5', label: '近3小时'},
        {key: '6', label: '近6小时'},
        {key: '7', label: '近12小时'},
        {key: '8', label: '近24小时'},

    ];
    return (
        <>
            <RangePicker showTime onChange={onTimeRangeChange}/>
            <Dropdown menu={{items, onClick}}>
                <Button>
                    {selectedOption || <DownOutlined/>}
                </Button>
            </Dropdown>        </>
    )
}


const SyncSelector: React.FC = () => {
    const [selectedOption, setSelectedOption] = useState<string | null>(null);
    const [timeRange, setTimeRange] = useState<[number, number] | null>(null);
    const [refresh, setRefresh] = useState<string>(null)
    const [from, setFrom] = useState();
    const [to, setTo] = useState();
    const items: MenuProps['items'] = [
        {key: '1', label: 'off'},
        {key: '2', label: '5s'},
        {key: '3', label: '10s'},
        {key: '4', label: '30s'},
        {key: '5', label: '1m'},
        {key: '6', label: '5m'},
        {key: '7', label: '15m'},
        {key: '8', label: '30m'},
        {key: '9', label: '1h'},
        {key: '10', label: '2h'},
        {key: '11', label: '1d'},
    ];

    const onClick: MenuProps['onClick'] = ({key, label}) => {
        if (key === '1') {
            setSelectedOption(null);
            setRefresh(null)
            return
        }
        const selectedItem = items.find((item) => item.key === key);
        if (selectedItem) {
            setSelectedOption(selectedItem.label as string);
            // 设置刷新间隔
            setRefresh(selectedItem.label as string);
        }

    };


    const manualRefresh = () => {
        message.info("点击")
    }

    const handleTimeRangeChange = (dates: [moment.Moment, moment.Moment] | null) => {
        if (dates) {
            const startTimestamp = dates[0].valueOf();
            const endTimestamp = dates[1].valueOf();
            setFrom(startTimestamp)
            setTo(endTimestamp)
            console.log("Selected Time Range (milliseconds):", [startTimestamp, endTimestamp]);
            setTimeRange([startTimestamp, endTimestamp]);
        } else {
            console.log("No time range selected");
        }
    };

    return (
        <>
            <Flex justify={"flex-end"}>
                {/*时间选择*/}
                <TimeSelector onTimeRangeChange={handleTimeRangeChange}></TimeSelector>
                {/*    刷新*/}
                <Button icon={<SyncOutlined/>} onClick={manualRefresh}></Button>
                {/*    自动刷新选择*/}
                <Dropdown menu={{items, onClick}}>
                    <Button>
                        {selectedOption || <DownOutlined/>}
                    </Button>
                </Dropdown>

            </Flex>
        </>
    )
}

export default SyncSelector;
