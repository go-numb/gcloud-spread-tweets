import { useState } from 'react'
import BtnSign from '/public/assets/images/sign-in-with-twitter-gray.png'
// import BtnSign from './assets/sign-in-with-twitter-gray.png'
import './App.css'

import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import 'react-tabs/style/react-tabs.css';

function App() {
  const [msg, setMsg] = useState("")
  const [result, setResult] = useState("")
  const [token, setToken] = useState("")
  const requestOAuth = async () => {
    const url = import.meta.env.VITE_API_URL + '/api/x/auth/request'
    const result = await fetch(url, {
      method: 'GET',
    })

    // エラーを表示する
    if (!result.ok) {
      const data = await result.json()
      console.log(data.code, data.message)
      return
    }

    const data = await result.json()
    
    if (data.data != null && data.data != "") {
      console.log(data.code, data.message,  data.data.url)
      setToken(data.data.token)
      window.location.href = data.data.url
    }
  }

  const downloadPDF = async () => {
    // onclickでPDFをダウンロードする
    const url = import.meta.env.BASE_URL + "/assets/files/managed_sheets.xlsx"
    console.log(url);
    const a = document.createElement("a");
    a.href = url;
    a.download = 'managed_sheets.xlsx';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
  }

  const registor = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    console.log(import.meta.env.VITE_API_URL);


    const url = `${import.meta.env.VITE_API_URL}/api/spreadsheet/upload?spreadsheet_id=${e.currentTarget.spreadsheet_id.value}&token=${token}`
    const result = await fetch(url, {
      method: 'GET',
    })

    const data = await result.json()
    if (data.data != null && data.data != "") {
      console.log(data.data);
      
      setResult(data.data.name + "を登録しました。認証したTwitter/XアカウントとSpreadsheetで登録したUserIDが同一である必要があります。")
    }
  }

  const handler = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMsg(e.target.value)
  }

  return (
    <>
      <Tabs>
        <TabList>
          <Tab>1:FIRST</Tab>
          <Tab>2:SECOND</Tab>
          <Tab>3:FINAL</Tab>
          <Tab>-:DESCRIPT</Tab>
        </TabList>

        <TabPanel>
          <div>
            <h2>Spreadsheetの設置</h2>
            <dl style={{ textAlign: "left" }}>
              <dt>1. Excelをダウンロード</dt>
              <dd>下記ダウンロードボタンをクリックしてExcelをダウンロードして任意の場所に保存してください。</dd>
              <dt>2. Google spreadsheetを活用</dt>
              <dd>Google spreadsheetまたはDriveに当ファイルを移動保存し、Spreadsheet型式に変換してください。</dd>
              <dt>3. ファイル共有</dt>
              <dd>ファイル共有からURL共有を選択し、リンクURLを取得してください。</dd>
            </dl>
            <a onClick={downloadPDF}>
              <button className="btn-dl">Download Spreadsheet</button>
            </a>
          </div>
        </TabPanel>
        <TabPanel>
          <div>
            <h2>Twitter認証</h2>
            <dl>
              <dt>1. 認証ボタンをクリック</dt>
              <dd>Twitter API OAuthを通じて、Twitter POST権限を取得します。</dd>
              <dt>2. Twitter認証ページで承認</dt>
              <dd>Twitterから当アプリに対して、投稿権限が与えられます。</dd>
              <dt>3. Twitter投稿を自動化</dt>
              <dd>当権限にてTwitter投稿を自動化します。続けて、SpreadsheetIDを登録ください。</dd>
            </dl>
            <a onClick={requestOAuth} style={{ cursor: "pointer" }}>
              <img src={BtnSign} alt="Sign in with Twitter/X" />
            </a>
          </div>
        </TabPanel>
        <TabPanel>
          <div>
            <h2>登録</h2>
            <dl>
              <dt>1. SpreadsheetIDを登録</dt>
              <dd>SpreadsheetIDを登録し、自動化を開始してください。</dd>
              <dt>2. 登録が完了しましたら、自動化を開始します。</dt>
              <dd>認証したTwitter/Xアカウントと登録したSpreadsheet AccountIDを照合し、Twitter投稿を自動化します。認証したアカウントと登録するアカウントが同様のものであることを確認してください。</dd>
            </dl>

            <div dangerouslySetInnerHTML={{ __html: result }}></div>

            <form onSubmit={(e) => registor(e)}>
              <input type="text" id="spreadsheet_id" name="spreadsheet_id" onChange={handler} placeholder="Spreadsheet ID" value={msg}></input>
              <input type="submit" value="登録" />
            </form>
          </div>
        </TabPanel>
        <TabPanel>
          <h2>説明文</h2>
        </TabPanel>
      </Tabs>
    </>
  )
}

export default App
