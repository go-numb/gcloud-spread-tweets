import React, { useState, useEffect } from 'react';
import { Divider, Row, Col, Steps } from 'antd';

// default component for step
const StepFlow: React.FC<{ stepN: number }> = ({ stepN }) => {
    const [current, setCurrent] = useState(1);
    const onChange = (value: number) => {
        setCurrent(value);
    };

    useEffect(() => {
        setCurrent(stepN);
    }, [stepN]);

    return (
        <Row>
            <Col>
                <Divider />
                <Steps
                    current={current}
                    onChange={onChange}
                    direction="vertical"
                    items={[
                        {
                            title: '説明文を読む',
                            description: 'まず、提供された説明文をよく読んでください。この説明文には、これから行う手順の概要と、必要な情報が記載されています。不明な点がある場合は、ご質問ください。',
                        },
                        {
                            title: 'Spreadsheetをダウンロード',
                            description: '説明文に記載されているリンクをクリックして、Spreadsheetをダウンロードしてください。ダウンロードしたファイルを開き、内容を確認してください。このSpreadsheetは、これから行う作業に必要なテンプレートになります。Spreadsheetを開いたら、認証に必要な情報を入力してください。',
                        },
                        {
                            title: 'Twitter/Xアカウントで認証',
                            description: 'Spreadsheetに記載した投稿を行うアカウントで、Twitter/X認証を行ってください。Twitter/Xアカウントのユーザー名とパスワードを入力し、認証プロセスを完了してください。',
                        },
                        {
                            title: 'Spreadsheetの登録',
                            description: '認証が完了したら、Spreadsheetの入力完了を確認し、SpreadsheetIDまたはリンクを登録ください。SpreadsheetはURL共有状態であることを確認してください。',
                        },
                        {
                            title: '投稿データの登録・編集',
                            description: 'Spreadsheetの登録が完了したら、任意で投稿データの確認・新規登録・編集を行ってください。新規投稿フォームから、投稿したい内容を入力してください。登録済みの投稿を編集・削除を行う際は登録済みの投稿を読み込無事で行なえます。編集入力が完了したら、表列末尾の編集ボタンを、削除であれば削除ボタンをクリックしてください。',
                        },
                        {
                            title: '作業の完了',
                            description: '以上の手順が完了したら、作業は終了です。そのままウィンドウまたはアプリを閉じ終了します。また、作業中に発生した問題やご質問がある場合は、遠慮なくお問い合わせください。',
                        },
                    ]}
                />
            </Col>
        </Row>
    );
}

export default StepFlow;