import React, {useEffect, useState} from 'react';
import {Space} from 'antd';
import SyncSelector from "@/pages/Metrics/components/SyncSelector";
import MetricsPanel from "@/pages/Metrics/components/MetricPanel";


const Metrics: React.FC = () => {
  const [timeRange, setTimeRange] = useState<[number, number] | null>(null);
  const [refresh, setRefresh] = useState<string | null>(null);
  const [refreshCount, setRefreshCount] = useState<number>(0);

  useEffect(() => {
    console.log("Refreshing Metrics due to refresh change or manual refresh.");
    // Increment refreshCount to force MetricsPanel to reload
    setRefreshCount((prevCount) => prevCount + 1);
  }, [refresh, timeRange]);

  const getGrafanaUrl = (panelId: string): string => {
    let url = `http://localhost:3000/d-solo/-lPyM7wIk/short-link?orgId=1&panelId=${panelId}`;
    // if (timeRange) {
    //   url += `&from=${timeRange[0]}&to=${timeRange[1]}`;
    // }
    if (timeRange) {
      url += `&from=now-5m&to=now`;
    }
    if (refresh) {
      url += `&refresh=${refresh}`;
    }
    console.log(url)
    return url;
  };

  return (
    <>
      <SyncSelector onRefreshChange={setRefresh} onTimeRangeChange={setTimeRange}/>
      <Space direction="vertical" size="middle" style={{display: 'flex'}}>
        <MetricsPanel key={`panel-6-${refreshCount}`} grafanaUrl={getGrafanaUrl("6")}/>
        <MetricsPanel key={`panel-4-${refreshCount}`} grafanaUrl={getGrafanaUrl("4")}/>
        <MetricsPanel key={`panel-2-${refreshCount}`} grafanaUrl={getGrafanaUrl("2")}/>
      </Space>
    </>
  )
}

export default Metrics;
