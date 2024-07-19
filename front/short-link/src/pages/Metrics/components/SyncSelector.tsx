import React, {useState} from 'react';
import type {MenuProps} from 'antd';
import {Button, DatePicker, Dropdown, Flex, message} from 'antd';
import {DownOutlined, SyncOutlined} from "@ant-design/icons";
import moment from 'moment';

//
const {RangePicker} = DatePicker;
//
//
// const TimeSelector: React.FC<{
//   onTimeRangeChange: (dates: [moment.Moment, moment.Moment] | null) => void
// }> = ({onTimeRangeChange}) => {
//
//   const [selectedOption, setSelectedOption] = useState<string | null>(null);
//   const onClick: MenuProps['onClick'] = ({key, label}) => {
//     const selectedItem = items.find((item) => item.key === key);
//     if (selectedItem) {
//       setSelectedOption(selectedItem.label as string);
//       onTimeRangeChange()
//     }
//     // 设置
//
//   };
//   const items: MenuProps['items'] = [
//     {key: '1', label: '近5分钟'},
//     {key: '2', label: '近15分钟'},
//     {key: '3', label: '近30分钟'},
//     {key: '4', label: '近1小时'},
//     {key: '5', label: '近3小时'},
//     {key: '6', label: '近6小时'},
//     {key: '7', label: '近12小时'},
//     {key: '8', label: '近24小时'},
//
//   ];
//   return (
//     <>
//       <RangePicker showTime onChange={onTimeRangeChange}/>
//       <Dropdown menu={{items, onClick}}>
//         <Button>
//           {selectedOption || <DownOutlined/>}
//         </Button>
//       </Dropdown>        </>
//   )
// }


const TimeSelector: React.FC<{
  onTimeRangeChange: (dates: [moment.Moment, moment.Moment] | null) => void
}> = ({onTimeRangeChange}) => {

  const [selectedOption, setSelectedOption] = useState<string | null>(null);

  const onClick: MenuProps['onClick'] = ({key}) => {
    let startMoment: moment.Moment = moment();
    let endMoment: moment.Moment = moment();

    switch (key) {
      case '1':
        startMoment = moment().subtract(5, 'minutes');
        break;
      case '2':
        startMoment = moment().subtract(15, 'minutes');
        break;
      case '3':
        startMoment = moment().subtract(30, 'minutes');
        break;
      case '4':
        startMoment = moment().subtract(1, 'hour');
        break;
      case '5':
        startMoment = moment().subtract(3, 'hours');
        break;
      case '6':
        startMoment = moment().subtract(6, 'hours');
        break;
      case '7':
        startMoment = moment().subtract(12, 'hours');
        break;
      case '8':
        startMoment = moment().subtract(24, 'hours');
        break;
      default:
        break;
    }

    setSelectedOption(key); // 设置选中的选项

    if (key === 'now') {
      onTimeRangeChange(null); // 点击 'now' 选项时，清空时间范围
    } else {
      onTimeRangeChange([startMoment, endMoment]); // 设置时间范围
    }
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
    {key: 'now', label: '现在'}, // 添加 '现在' 选项
  ];

  return (
    <>
      <RangePicker showTime onChange={onTimeRangeChange}/>
      <Dropdown menu={{items, onClick}}>
        <Button>
          {selectedOption ? items.find(item => item.key === selectedOption)?.label || <DownOutlined/> : <DownOutlined/>}
        </Button>
      </Dropdown>
    </>
  )
}

interface SyncSelectorProps {
  onRefreshChange: (refresh: string | null) => void;
  onTimeRangeChange: (timeRange: [number, number] | null) => void;
}

const SyncSelector: React.FC<SyncSelectorProps> = ({onRefreshChange, onTimeRangeChange}) => {
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
      // 传递给父亲组件
      onRefreshChange(selectedItem.label as string);

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
      onTimeRangeChange([startTimestamp, endTimestamp]);
    } else {
      console.log("No time range selected");
      onTimeRangeChange(null);
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
