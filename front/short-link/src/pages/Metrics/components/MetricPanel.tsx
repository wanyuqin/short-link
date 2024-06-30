import React from 'react';
import { Col, Row } from 'antd';

interface MetricsPanelProps {
    grafanaUrl: string;
}

const MetricsPanel: React.FC<MetricsPanelProps> = ({ grafanaUrl }) => {
    return (
        <Row>
            <Col span={24}>
                <iframe
                    src={grafanaUrl}
                    width="100%" height="400" frameBorder="0"
                    style={{ border: 'none', width: '100%', minHeight: '400px' }}
                ></iframe>
            </Col>
        </Row>
    );
}

export default MetricsPanel;
