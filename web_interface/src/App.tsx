import { useEffect, useState } from 'react'
import { useLocation } from 'react-router-dom';
import BtnSign from '/public/assets/images/sign-in-with-twitter-gray.png'
// import BtnSign from './assets/sign-in-with-twitter-gray.png'
import './App.css'

import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import 'react-tabs/style/react-tabs.css';


function App() {
  const location = useLocation();
  const [msg, setMsg] = useState("")
  const [result, setResult] = useState("")
  const [token, setToken] = useState("")
  const [username, setUsername] = useState("")

  const requestOAuth = async () => {
    const url = import.meta.env.VITE_API_URL + '/api/x/auth/request'
    const res = await fetch(url, {
      method: 'GET',
    })

    // エラーを表示する
    if (!res.ok) {
      const data = await res.json()
      console.debug(data.code, data.message)
      return
    }

    const data = await res.json()

    if (data.data != null && data.data != "") {
      console.debug(data.code, data.message, data.data.url)
      window.location.href = data.data.url
    }
  }

  const downloadPDF = async () => {
    // onclickでPDFをダウンロードする
    const url = import.meta.env.BASE_URL + "/assets/files/managed_sheets.xlsx"
    console.debug(url);
    try {
      const a = document.createElement("a");
      a.href = url;
      a.download = 'managed_sheets.xlsx';
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
    } catch (err) {
      console.error(err);
      setResult(`<b>ダウンロードエラーが発生しました。<br />${err}</b>`)
    }
  }

  const registor = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    console.debug(import.meta.env.VITE_API_URL);


    const url = `${import.meta.env.VITE_API_URL}/api/spreadsheet/upload?spreadsheet_id=${e.currentTarget.spreadsheet_id.value}&token=${token}`
    const result = await fetch(url, {
      method: 'GET',
    })

    const data = await result.json()
    // エラーを表示する
    if (!result.ok) {
      setMsg("")
      setResult(`<b>エラーが発生しました。もう一度お試しください。<br />${data.code}: ${data.message}</b>`)
      return
    } else if (data.data != null && data.data != "") {
      console.debug(data.data);

      setMsg("")
      setResult("<b>Twitter/X Account: [ " + data.data.name + " ] を登録しました。明日より、自動投稿を行います。認証したTwitter/Xアカウントと登録アカウントが一致していることを確認してください。</b>")
    }
  }

  const handler = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMsg(e.target.value)
  }

  // urlからtokenを取得
  useEffect(() => {
    const query = new URLSearchParams(location.search);
    const token = query.get('token');
    if (token != null) {
      setToken(token);
    }
    const username = query.get('username');
    if (username != null) {
      setUsername(username);
      setResult(`<b>Twitter/X Account: [ ${username} ] で認証を得ました。次に、SpreadsheetIDを登録してください。</b>`)
    }
  }, [location]);

  const authAccount = () => {
    if (username != "") {
      return `<p><b>認証済みTwitter/X Account: [ ${username} ]</b></p>`
    }
    return "Twitter/X Account: [ 未認証 ]"
  };

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
              <dd>認証したTwitter/Xアカウント[ {username} ]と登録したSpreadsheet AccountIDを照合し、Twitter投稿を自動化します。認証したアカウントと登録するアカウントが同様のものであることを確認してください。
              <div dangerouslySetInnerHTML={{ __html: authAccount() }}></div>
              </dd>
            </dl>

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

      <div style={{display: "block", position: "absolute", left: 0, right: 0, bottom: 0, justifyContent: "center", alignItems: "center"}}>
      <hr />
        <div dangerouslySetInnerHTML={{ __html: result }}></div>
      </div>
    </>
  )
}

export default App
