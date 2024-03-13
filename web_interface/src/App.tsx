import { useEffect, useState } from 'react'
import { useLocation } from 'react-router-dom';
import BtnSign from '/public/assets/images/sign-in-with-twitter-gray.png'
// import BtnSign from './assets/sign-in-with-twitter-gray.png'
import './App.css'
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { Button, Flex, Layout, Col, Row, Input, Popconfirm, message, Spin } from 'antd';
const { Footer, Content } = Layout;
const { TextArea } = Input;


import Indicator from './components/data';
import StepFlow from './components/step';

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

// const headerStyle: React.CSSProperties = {
// };

const contentStyle: React.CSSProperties = {
  textAlign: 'left',
};

const footerStyle: React.CSSProperties = {
  backgroundColor: '#fcfcfc',
  textAlign: 'center',
};

const layoutStyle = {
  borderRadius: 8,
  overflow: 'hidden',
  width: 'calc(50% - 8px)',
  maxWidth: 'calc(50% - 8px)',
};


function App() {
  const location = useLocation();
  const [spreadsheetId, setSpreadsheetId] = useState("")
  const [result, setResult] = useState("")
  const [token, setToken] = useState("")
  const [username, setUsername] = useState("")
  const [spreadsheetIds, setSpreadsheetIds] = useState<string[]>([])

  // current stepN state
  const [current, setCurrent] = useState(0);

  // Loading state
  const [loading, setLoading] = useState<boolean>(false);

  // antdのmessageを使用する
  const [messageApi, contextHolder] = message.useMessage();

  // formから得た値を元に、postを登録する
  const createPost = (e: React.FormEvent<HTMLFormElement>) => {
    setLoading(true);

    e.preventDefault()
    const postText = e.currentTarget.postText.value
    const file1 = e.currentTarget.file1.value
    const file2 = e.currentTarget.file2.value
    const file3 = e.currentTarget.file3.value
    const file4 = e.currentTarget.file4.value
    const priority = parseInt(e.currentTarget.priority.value, 10)

    if (postText === "") {
      messageApi.error("失敗しました、投稿本文が空です。")
      return
    }

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
        messageApi.success(`post id: [ ${res.data.data.uuid} ] の投稿を登録しました。`);
        setResult(`<b>post id: [ ${res.data.data.uuid} ] の投稿を登録しました。</b>`)
      })
      .catch((err) => {
        console.error(err)
        messageApi.error(`エラーが発生しました。もう一度お試しください。\n${err}`);
        setResult(`<b>エラーが発生しました。もう一度お試しください。<br />${err}</b>`)
      })
      .finally(() => {
        setLoading(false);
        setCurrent(4)
      });
  }

  // posts:Post[]からuuidで検索してPostを取得する
  const putPost = (uuid: String) => {
    setLoading(true);

    console.debug(uuid);
    const send_post = posts.find((post) => post.uuid === uuid)
    console.debug(send_post)

    // popupで確認を行い、更新を行う
    console.debug("更新します", send_post)
    const url = import.meta.env.VITE_API_URL + '/api/x/post?token=' + token + '&username=' + username
    axios.put(url, send_post)
      .then((res) => {
        messageApi.success(`post id: [ ${res.data.data.uuid} ] の投稿を更新しました。`);
        setResult(`<b>post id: [ ${res.data.data.uuid} ] の投稿を更新しました。</b>`)
      })
      .catch((err) => {
        console.error(err)
        messageApi.error(`エラーが発生しました。もう一度お試しください。\n${err}`);
        setResult(`<b>エラーが発生しました。もう一度お試しください。<br />${err}</b>`)
      })
      .finally(() => {
        setLoading(false);
        setCurrent(5)
      });
  }

  // posts:Post[]からuuidで検索してPostを取得する
  const deletePost = (uuid: String) => {
    setLoading(true);

    console.debug(uuid);
    // popupで確認を行い、更新を行う
    console.debug("削除します")
    const url = import.meta.env.VITE_API_URL + '/api/x/post?token=' + token + '&username=' + username + '&uuid=' + uuid
    axios.delete(url)
      .then((res) => {
        // postsからuuidを持つpostを削除
        const newPosts = posts.filter((post) => post.uuid !== uuid)
        setPosts(newPosts)
        messageApi.success(`post id: [ ${res.data.data} ] の投稿を削除しました。`);
        setResult(`<b>post id: [ ${res.data.data} ] の投稿を削除しました。</b>`)
      })
      .catch((err) => {
        console.error(err)
        messageApi.error(`エラーが発生しました。もう一度お試しください。\n${err}`);
        setResult(`<b>エラーが発生しました。もう一度お試しください。<br />${err}</b>`)
      })
      .finally(() => {
        setLoading(false)
        setCurrent(5)
      });
  }

  const requestOAuth = async () => {
    setLoading(true);

    const url = import.meta.env.VITE_API_URL + '/api/x/auth/request'
    axios.get(url)
      .then(res => {
        console.log(res.data);
        window.location.href = res.data.data.url;
      })
      .catch(error => {
        console.log(error);
        const message = `<b>エラーが発生しました。もう一度お試しください。<br />${error}</b>`;
        messageApi.error(message);
        setResult(message)
      })
      .finally(() => {
        setLoading(false);
        setCurrent(3)
      });
  }

  const downloadPDF = async () => {
    setLoading(true);

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
      const message = `<b>ダウンロードエラーが発生しました。<br />${err}</b>`;
      messageApi.error(message);
      setResult(message)
    } finally {
      setLoading(false);
      setCurrent(2)
    }
  }

  const registor = async (e: React.FormEvent<HTMLFormElement>, send_type: string) => {
    setLoading(true);

    e.preventDefault()

    console.debug(import.meta.env.VITE_API_URL);

    let url = `${import.meta.env.VITE_API_URL}/api/spreadsheet/upload?spreadsheet_id=${e.currentTarget.spreadsheet_id.value}&token=${token}`
    if (send_type === "repost") {
      url = url + "&repost=true"
    }

    axios.get(url)
      .then(res => {
        setSpreadsheetId("")
        messageApi.success("Spreadsheet ID: [ " + res.data.data.spreadsheet_id + " ] を取得し、Twitter/X Account: [ " + res.data.data.id + " ] と合致する投稿を登録しました。明日より、自動投稿を行います。");
        setResult("<b>Spreadsheet ID: [ " + res.data.data.spreadsheet_id + " ] を取得し、Twitter/X Account: [ " + res.data.data.id + " ] と合致する投稿を登録しました。明日より、自動投稿を行います。</b>")
        setSpreadsheetIds((rev) => [...rev, res.data.data.spreadsheet_id])
      })
      .catch(err => {
        console.log(err);
        messageApi.error(`エラーが発生しました。もう一度お試しください。\n${err}`);
        setResult(`<b>エラーが発生しました。もう一度お試しください。<br />${err.code}: ${err.message}</b>`)
      })
      .finally(() => {
        setLoading(false);
        setCurrent(4)
      });
  }


  const [posts, setPosts] = useState<Post[]>([])
  const getPost = async () => {
    setLoading(true);

    const url = import.meta.env.VITE_API_URL + '/api/x/posts?token=' + token + "&username=" + username
    axios.get(url)
      .then(res => {
        // json dataをPost型に変換
        const tempPosts: Post[] = res.data.data
        setPosts(tempPosts)

        console.debug(res.data.code, res.data.message, res.data.data)
        messageApi.success(`UserID: [ ${username} ] の投稿を取得しました。`);
        setResult(`<b>UserID: [ ${username} ] の投稿を取得しました。</b>`)
      })
      .catch(err => {
        console.debug(err.code, err.message)
        messageApi.error(`エラーが発生しました。もう一度お試しください。\n${err}`);
      })
      .finally(() => {
        setLoading(false);
        setCurrent(4)
      });
  }


  const handler = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSpreadsheetId(e.target.value)
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

  const handleTextAreaChange = (event: React.ChangeEvent<HTMLTextAreaElement>, index: number, field: string) => {
    const newPosts: Post[] = [...posts];

    (newPosts[index] as any)[field] = event.target.value;

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
      messageApi.success("Hello, " + username + "!");
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
      <div style={{ display: "flex", justifyContent: "center", alignItems: "center" }}>
        {contextHolder}
      </div>

      <div id="root">
        <Flex gap='middle' wrap='wrap'>
          <Layout style={layoutStyle}>
            {/* <Header style={headerStyle}>

        </Header> */}
            <Content style={contentStyle}>
              <Tabs>
                <Spin spinning={loading} delay={0}>
                  <TabList>
                    <Tab>0:説明</Tab>
                    <Tab>1:FIRST</Tab>
                    <Tab>2:SECOND</Tab>
                    <Tab>3:FINAL</Tab>
                    <Tab>4:USERS</Tab>
                  </TabList>
                </Spin>

                <TabPanel>
                  <Row style={{ paddingBottom: "3rem" }}>
                    <Col>
                    {/* Indicatorに引数を渡す */}
                      <Indicator />
                      <h2>説明文</h2>
                      <StepFlow stepN={current} />
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
                    </Col>
                  </Row>
                </TabPanel>
                <TabPanel>
                  <Row style={{ paddingBottom: "3rem" }}>
                    <Col>
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
                        <Button type="link">Download Spreadsheet</Button>
                      </a>

                    </Col>
                  </Row>
                </TabPanel>
                <TabPanel>
                  <Row>
                    <Col>
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
                    </Col>
                  </Row>
                </TabPanel>
                <TabPanel>
                  <Row>
                    <Col>
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
                        <Row>
                          <Col span={6}>
                            <Input type="text" maxLength={100} id="spreadsheet_id" name="spreadsheet_id" onChange={handler} placeholder="Spreadsheet ID" value={spreadsheetId} />
                          </Col>
                          <Col>
                            <Input type="submit" value="登録" />
                          </Col>
                        </Row>
                      </form>
                    </Col>
                  </Row>
                </TabPanel>
                <TabPanel>
                  <Row>
                    <Col>
                      <h2>データ操作</h2>
                      <p>
                        <b>Account ID: [ {username} ] の投稿を取得/表示します。</b>
                      </p>
                      <Row className='m-2rem'>
                        <Col>
                          <Button type='primary' onClick={getPost}>Get posts</Button>
                        </Col>
                      </Row>


                      {/* postsをテーブル型式のリストで表示する */}
                      <div className="container">
                        <table>
                          <thead>
                            <tr>
                              <th>ID</th>
                              <th>投稿本文</th>
                              <th>ファイル1</th>
                              <th>ファイル2</th>
                              <th>ファイル3</th>
                              <th>ファイル4</th>
                              <th>ファイル必須</th>
                              <th>投稿可否</th>
                              <th>優先度</th>
                              <th>操作</th>
                            </tr>
                          </thead>
                          <tbody>
                            {posts.map((post, index) => (
                              <tr key={post.uuid}>
                                <td>{post.id}</td>
                                <td><TextArea cols={100} value={post.text} onChange={(e) => handleTextAreaChange(e, index, 'text')} /></td>
                                <td><Input type="text" value={post.file_1} onChange={(e) => handleInputChange(e, index, 'file_1')} /></td>
                                <td><Input type="text" value={post.file_2} onChange={(e) => handleInputChange(e, index, 'file_2')} /></td>
                                <td><Input type="text" value={post.file_3} onChange={(e) => handleInputChange(e, index, 'file_3')} /></td>
                                <td><Input type="text" value={post.file_4} onChange={(e) => handleInputChange(e, index, 'file_4')} /></td>
                                <td><Input type="number" value={post.with_files} onChange={(e) => handleInputChange(e, index, 'with_files')} /></td>
                                <td><Input type="number" value={post.checked} onChange={(e) => handleInputChange(e, index, 'checked')} /></td>
                                <td><Input type="number" value={post.priority} onChange={(e) => handleInputChange(e, index, 'priority')} /></td>
                                <td>
                                  <Row>
                                    <Col>
                                      <Popconfirm
                                        title="投稿の更新"
                                        description="投稿の更新を行います。"
                                        onConfirm={() => putPost(post.uuid)}
                                        // onCancel={cancel}
                                        okText="ok"
                                        cancelText="cancel"
                                      >
                                        <Button title='更新'>📝</Button>
                                      </Popconfirm>
                                    </Col>
                                  </Row>
                                  <Row>
                                    <Col>
                                      <Popconfirm
                                        title="投稿の削除"
                                        description="投稿の削除を行います。"
                                        onConfirm={() => deletePost(post.uuid)}
                                        icon={<ExclamationCircleOutlined style={{ color: 'red' }} />}
                                        // onCancel={cancel}
                                        okText="ok"
                                        cancelText="cancel"
                                      >
                                        <Button title='削除'>🗑</Button>
                                      </Popconfirm>

                                    </Col>
                                  </Row>
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
                          <Input type="text" id="id" name="id" value={username} disabled />
                          <label htmlFor="postText">Post Text:</label>
                          <TextArea style={{ margin: "1rem" }} id="postText" name="postText" placeholder="Post Text" showCount maxLength={25000} />
                          <label htmlFor="file1">File 1:</label>
                          <Input type="text" id="file1" name="file1" placeholder="File 1" />
                          <label htmlFor="file2">File 2:</label>
                          <Input type="text" id="file2" name="file2" placeholder="File 2" />
                          <label htmlFor="file3">File 3:</label>
                          <Input type="text" id="file3" name="file3" placeholder="File 3" />
                          <label htmlFor="file4">File 4:</label>
                          <Input type="text" id="file4" name="file4" placeholder="File 4" />
                          <label htmlFor="priority">Priority:</label>
                          <Input type="number" step={1} id="priority" name="priority" placeholder="Priority" />
                          <Input type="submit" value="新規登録" />
                        </form>
                      </div>


                      <div style={{ margin: "3rem auto", paddingBottom: "5rem" }}>
                        <h2>SpreadsheetIDを再登録</h2>
                        <form onSubmit={(e) => registor(e, "repost")}>
                          <Row>
                            <Col span={6}>
                              <Input type="text" id="spreadsheet_id" name="spreadsheet_id" onChange={handler} placeholder="Spreadsheet ID" value={spreadsheetId} />
                            </Col>
                            <Col>
                              <Input type="submit" value="登録" />
                            </Col>
                          </Row>
                        </form>
                      </div>
                    </Col>
                  </Row>
                </TabPanel>
              </Tabs>
            </Content>
          </Layout>

          <Footer style={footerStyle} className='fixed-bottom'>
            <hr />
            <div dangerouslySetInnerHTML={{ __html: result }}></div>
          </Footer>
        </Flex>
      </div>
    </>
  )
}

export default App
