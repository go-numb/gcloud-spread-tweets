import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { ArrowDownOutlined, ArrowUpOutlined } from '@ant-design/icons';
import { Card, Col, Row, Statistic } from 'antd';

const UPCOLOR = '#3f8600';
const DOWNCOLOR = '#cf1322';


const Indicator: React.FC = () => {
    const [userN, setUserN] = useState<number>(0);
    const [postN, setPostN] = useState<number>(0);
    const [latestUserN, setLatestUserN] = useState<number>(0);
    const [latestPostN, setLatestPostN] = useState<number>(0);

    const getData = () => {
        const url = import.meta.env.VITE_API_URL + '/api/data/usage';
        console.log(url);

        axios.get(url)
            .then(res => {
                console.log(res.data);
                setUserN(res.data.data.users);
                setPostN(res.data.data.posts);
            })
            .catch(error => {
                console.log(error);
            });
    }

    // 呼び出し時、一度だけ実行する関数
    useEffect(() => {
        getData();
        if (userN === 0) {
            setLatestUserN(userN);
            setLatestPostN(postN);
        }
    }, []);


    return (
        <>
            <Row gutter={16}>
                <Col span={12}>
                    <Card bordered={false}>
                        <Statistic
                            title="Active accounts"
                            value={userN}
                            precision={0}
                            valueStyle={userN >= latestUserN ? { color: UPCOLOR } : { color: DOWNCOLOR }}
                            // userNがlatestUserNより大きい場合、上向きの矢印を表示
                            prefix={userN >= latestUserN ? <ArrowUpOutlined /> : <ArrowDownOutlined />}
                            suffix="accounts"
                        />
                    </Card>
                </Col>
                <Col span={12}>
                    <Card bordered={false}>
                        <Statistic
                            title="Active posts"
                            value={postN}
                            precision={0}
                            valueStyle={postN >= latestPostN ? { color: UPCOLOR } : { color: DOWNCOLOR }}
                            // postNがlatestPostNより大きい場合、上向きの矢印を表示
                            prefix={postN >= latestPostN ? <ArrowUpOutlined /> : <ArrowDownOutlined />}
                            suffix="posts"
                        />
                    </Card>
                </Col>
            </Row>
        </>
    );
};

export default Indicator;