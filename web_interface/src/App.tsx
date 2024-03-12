import { useEffect, useState } from 'react'
import { useLocation } from 'react-router-dom';
import BtnSign from '/public/assets/images/sign-in-with-twitter-gray.png'
// import BtnSign from './assets/sign-in-with-twitter-gray.png'
import './App.css'

import axios from 'axios';

import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import 'react-tabs/style/react-tabs.css';

interface Post {
  uuid: string; // "omitempty" が指定されているのでオプショナル
  id?: string;
  text?: string;
  file_1?: string;
  file_2?: string;
  file_3?: string;
  file_4?: string;
  with_files?: number;
  checked?: number;
  priority?: number;
  count?: number;
  post_url?: string;
  last_posted_at?: Date; // time.Time は Date 型に相当する
  created_at?: Date;
}


function App() {
  const location = useLocation();
  const [msg, setMsg] = useState("")
  const [result, setResult] = useState("")
  const [token, setToken] = useState("")
  const [username, setUsername] = useState("")
  const [spreadsheetIds, setSpreadsheetIds] = useState<string[]>([])

  // formから得た値を元に、postを登録する
  const createPost = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    const postText = e.currentTarget.postText.value
    const file1 = e.currentTarget.file1.value
    const file2 = e.currentTarget.file2.value
    const file3 = e.currentTarget.file3.value
    const file4 = e.currentTarget.file4.value
    const priority = parseInt(e.currentTarget.priority.value, 10)

    const url = import.meta.env.VITE_API_URL + '/api/x/post?token=' + token + '&username=' + username
    const send_post: Post = {
      uuid: "",
      id: username,
      text: postText,
      file_1: file1,
      file_2: file2,
      file_3: file3,
      file_4: file4,
      with_files: 1,
      checked: 1,
      priority: priority,
    }
    console.debug(send_post);

    axios.post(url, send_post)
      .then((res) => {
        // postsに新しいpostを追加
        const newPosts = [...posts, res.data.data]
        setPosts(newPosts)

        setResult(`<b>post id: [ ${res.data.data.uuid} ] の投稿を登録しました。</b>`)
      })
      .catch((err) => {
        console.error(err)
        setResult(`<b>エラーが発生しました。もう一度お試しください。<br />${err}</b>`)
      })
  }

  // posts:Post[]からuuidで検索してPostを取得する
  const putPost = (uuid: String) => {
    console.debug(uuid);
    const send_post = posts.find((post) => post.uuid === uuid)
    console.debug(send_post)

    // popupで確認を行い、更新を行う
    const result = window.confirm("uuid: " + uuid + "を更新します。")
    if (result) {
      console.debug("更新します", send_post)
      const url = import.meta.env.VITE_API_URL + '/api/x/post?token=' + token + '&username=' + username
      axios.put(url, send_post)
        .then((res) => {
          setResult(`<b>post id: [ ${res.data.data.uuid} ] の投稿を更新しました。</b>`)
        })
        .catch((err) => {
          console.error(err)
          setResult(`<b>エラーが発生しました。もう一度お試しください。<br />${err}</b>`)
        })

    } else {
      console.debug("更新しません")
      setResult(`<b>更新を取りやめました。</b>`)
    }
  }

  // posts:Post[]からuuidで検索してPostを取得する
  const deletePost = (uuid: String) => {
    console.debug(uuid);

    // popupで確認を行い、更新を行う
    const result = window.confirm("uuid: " + uuid + "を削除します。")
    if (result) {
      console.debug("削除します")
      const url = import.meta.env.VITE_API_URL + '/api/x/post?token=' + token + '&username=' + username + '&uuid=' + uuid
      axios.delete(url)
        .then((res) => {
          // postsからuuidを持つpostを削除
          const newPosts = posts.filter((post) => post.uuid !== uuid)
          setPosts(newPosts)
          setResult(`<b>post id: [ ${res.data.data} ] の投稿を削除しました。</b>`)
        })
        .catch((err) => {
          console.error(err)
          setResult(`<b>エラーが発生しました。もう一度お試しください。<br />${err}</b>`)
        })

    } else {
      console.debug("削除しません")
      setResult(`<b>削除を取りやめました。</b>`)
    }
  }


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

  const registor = async (e: React.FormEvent<HTMLFormElement>, send_type: string) => {
    e.preventDefault()

    console.debug(import.meta.env.VITE_API_URL);


    let url = `${import.meta.env.VITE_API_URL}/api/spreadsheet/upload?spreadsheet_id=${e.currentTarget.spreadsheet_id.value}&token=${token}`
    if (send_type === "repost") {
      url = url + "&repost=true"
    }
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
      setResult("<b>Spreadsheet ID: [ " + data.data.spreadsheet_id + " ] を取得し、Twitter/X Account: [ " + data.data.id + " ] と合致する投稿を登録しました。明日より、自動投稿を行います。</b>")
      setSpreadsheetIds((rev) => [...rev, data.data.spreadsheet_id])
    }
  }


  const [posts, setPosts] = useState<Post[]>([])
  const getPost = async () => {
    const url = import.meta.env.VITE_API_URL + '/api/x/posts?token=' + token + "&username=" + username
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
      // json dataをPost型に変換
      const tempPosts: Post[] = data.data
      setPosts(tempPosts)

      console.debug(data.code, data.message, data.data)
      setResult(`<b>UserID: [ ${username} ] の投稿を取得しました。</b>`)
    }
  }


  const handler = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMsg(e.target.value)
  }

  const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>, index: number, field: string) => {
    const newPosts: Post[] = [...posts];

    // Form値であるString型を数値型に変換
    if (field === 'priority' || field === 'with_files' || field === 'checked' || field === 'count') {
      // priorityとwith_files, checked, countフィールドは数値型に変換
      (newPosts[index] as any)[field] = Number(event.target.value);
    } else {
      // その他のフィールドは文字列型のまま
      (newPosts[index] as any)[field] = event.target.value;
    }

    setPosts(newPosts);
  };

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
      <div style={{ display: "block", padding: "5rem auto" }}>
        <Tabs>
          <TabList>
            <Tab>0:説明</Tab>
            <Tab>1:FIRST</Tab>
            <Tab>2:SECOND</Tab>
            <Tab>3:FINAL</Tab>
            <Tab>4:USERS</Tab>

          </TabList>

          <TabPanel>
            <div>
              <h2>説明文</h2>
              <h2># spread to twitter/x</h2>
              <ul>
                <li>Spreadsheetに記載されたアカウントと同じTwitter/Xアカウントで認証を行ってください。</li>
                <li>行列の構造を変更せず、x_text以外は半角で記載ください。</li>
                <li>yes/noを問う列は、1ならYes, 1以外ならばNoと判断されます。</li>
                <li>sheet: x_postsを埋めてください。</li>
                <li>x_files各セルは、Google Drive上のファイルかつURL共有されたURLまたはFileIDを記載ください。</li>
                <li>with_filesは、添付ファイルの投稿が失敗した場合、投稿を行うか、行わないかを0/1で記載ください。</li>
                <li>checkedは、投稿の可否であり、0/1で記載ください。</li>
              </ul>
              <br />
              記入を終え、Google Spreadsheetに保存し、Spreadsheetファイル形式・URL共有で保存しましたら、SpreadsheetIDまたは共有URLをブラウザから登録してください。
              <br />
              <br />
              <h2>有料版について</h2>
              <ul>
                <li>無料版はSpreadsheetを登録時1回、有料版は投稿時に都度読み込みます。</li>
                <li>割り振られた回数/月、5分おきの自動投稿が行なえます。</li>
                <li>時間指定を複数登録可能です。</li>
                <li>X Blue登録アカウントの場合、長文投稿が可能です.</li>
              </ul>
            </div>
          </TabPanel>
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
                <dd>認証したTwitter/Xアカウント[ {username} ]と登録したSpreadsheet AccountIDを照合し、Twitter投稿を自動化します。認証したアカウントと登録するアカウントが同様のものであることを確認してください.
                  <div dangerouslySetInnerHTML={{ __html: authAccount() }}></div>
                </dd>
                <dt>登録済みのSpreadsheetIDs</dt>
                {/* 配列を人間にわかりやすい適切な表記で表示 */}
                <dd>{spreadsheetIds.length == 0 ? "登録無し" : `[${spreadsheetIds.join(", ")}]`}</dd>
              </dl>

              <form onSubmit={(e) => registor(e, "")}>
                <input type="text" id="spreadsheet_id" name="spreadsheet_id" onChange={handler} placeholder="Spreadsheet ID" value={msg}></input>
                <input type="submit" value="登録" />
              </form>

            </div>
          </TabPanel>
          <TabPanel>
            <div>
              <h2>データ操作</h2>
              <p>
                <b>Account ID: [ {username} ] の投稿を取得/表示します。</b>
              </p>
              <input className='btn-dl' type="button" onClick={getPost} value="get" />
              {/* postsをテーブル型式のリストで表示する */}
              <div className="container">
                <table>
                  <thead>
                    <tr>
                      <th>AccountID</th>
                      <th>Text</th>
                      <th>File1</th>
                      <th>File2</th>
                      <th>File3</th>
                      <th>File4</th>
                      <th>WithFiles</th>
                      <th>Checked</th>
                      <th>Priority</th>
                      <th>Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {posts.map((post, index) => (
                      <tr key={post.uuid}>
                        <td>{post.id}</td>
                        <td><input type="text" value={post.text} onChange={(e) => handleInputChange(e, index, 'text')} /></td>
                        <td><input type="text" value={post.file_1} onChange={(e) => handleInputChange(e, index, 'file_1')} /></td>
                        <td><input type="text" value={post.file_2} onChange={(e) => handleInputChange(e, index, 'file_2')} /></td>
                        <td><input type="text" value={post.file_3} onChange={(e) => handleInputChange(e, index, 'file_3')} /></td>
                        <td><input type="text" value={post.file_4} onChange={(e) => handleInputChange(e, index, 'file_4')} /></td>
                        <td><input type="number" value={post.with_files} onChange={(e) => handleInputChange(e, index, 'with_files')} /></td>
                        <td><input type="number" value={post.checked} onChange={(e) => handleInputChange(e, index, 'checked')} /></td>
                        <td><input type="number" value={post.priority} onChange={(e) => handleInputChange(e, index, 'priority')} /></td>
                        <td>
                          <button onClick={() => deletePost(post.uuid)} title='削除'>🗑</button>
                          <button onClick={() => putPost(post.uuid)} title='更新'>📝</button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div style={{ margin: "3rem auto", paddingBottom: "5rem" }}>
                <h2>新規投稿追加フォーム</h2>
                <form className='form-horizon' onSubmit={(e) => createPost(e)}>
                  <label htmlFor="id">AccountID:</label>
                  <input type="text" id="id" name="id" value={username} disabled />
                  <label htmlFor="postText">Post Text:</label>
                  <textarea id="postText" name="postText" rows={10} cols={50} placeholder="Post Text" />
                  <label htmlFor="file1">File 1:</label>
                  <input type="text" id="file1" name="file1" placeholder="File 1" />
                  <label htmlFor="file2">File 2:</label>
                  <input type="text" id="file2" name="file2" placeholder="File 2" />
                  <label htmlFor="file3">File 3:</label>
                  <input type="text" id="file3" name="file3" placeholder="File 3" />
                  <label htmlFor="file4">File 4:</label>
                  <input type="text" id="file4" name="file4" placeholder="File 4" />
                  <label htmlFor="priority">Priority:</label>
                  <input type="number" step={1} id="priority" name="priority" placeholder="Priority" />
                  <input className='btn-dl' type="submit" value="Post" />
                </form>
              </div>


              <div style={{ margin: "3rem auto", paddingBottom: "5rem" }}>
                <h2>SpreadsheetIDを再登録</h2>
                <form onSubmit={(e) => registor(e, "repost")}>
                  <input type="text" id="spreadsheet_id" name="spreadsheet_id" onChange={handler} placeholder="Spreadsheet ID" value={msg}></input>
                  <input type="submit" value="登録" />
                </form>
              </div>
            </div>
          </TabPanel>
        </Tabs>
      </div>

      <div className='fixed-bottom'>
        <hr />
        <div dangerouslySetInnerHTML={{ __html: result }}></div>
      </div>
    </>
  )
}

export default App
