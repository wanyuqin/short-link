import React from 'react';
import {Space} from 'antd';
import SyncSelector from "@/pages/Metrics/components/SyncSelector";
import MetricsPanel from "@/pages/Metrics/components/MetricPanel";


const Metrics: React.FC = () => {
    const getGrafanaUrl = (panelId: string): string => {
        return `http://localhost:3000/d-solo/-lPyM7wIk/short-link?orgId=1&refresh=5s&panelId=${panelId}: ''}`;
    };
    return (
        <>
            <SyncSelector></SyncSelector>
            <Space direction="vertical" size="middle" style={{display: 'flex'}}>
                <MetricsPanel grafanaUrl={getGrafanaUrl("6")}/>
                <MetricsPanel grafanaUrl={getGrafanaUrl("4")}/>
                <MetricsPanel grafanaUrl={getGrafanaUrl("2")}/>
            </Space>
        </>
    )

}

export default Metrics;
