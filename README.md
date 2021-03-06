<h1>システム</h1>
<img src="./react/react-sample/src/images/logo.png">
<br/>
<br/>
<h1>インフラ構成図</h1>
<img src="./react/react-sample/src/images/manamana_arcitect.png">
<br/>
<br/>
<h1>使用技術</h1>
<p>◯Go 1.15.14</p>
<p>◯React 17.0.2</p>
<p>◯nginx 1.21.4</p>
<p>◯Dcoker 20.10.10<p>
<p>◯Terraform v1.1.3<p>
<p>◯GitHub</p>
<p>◯Visual Studio Code</p>
<br/>
<br/>
<h1>①どのような機能を</h1>
【機能概要】
<br/>
ipo（新規公開株式）の自動申し込み、<br/>
および現資産の確認・管理を行うシステムです。<br/>
忘れがちなipoの申し込みと、<br/>
資産運用の改善をこのアプリケーションで一括して行う事ができます。
<br/>
<br/>
【課題とその背景】<br/>
課題として、<br/>
◯SPA<br/>
◯インフラのコード化<br/>
◯GoとReactの使用<br/>
を実現したいと思っていました。
<br/>
<br/>
背景として、<br/>
日頃からこの様なシステムが欲しいと思っていましたが、<br/>
代替となるものは見つける事が出来なかったです。<br/>
また、今興味のある技術の学習を考えた時に、<br/>
自分が欲しいと思うシステムを開発する事がモチベーションに繋がると考え、<br/>
このシステムを作るに至りました。
<br/>
<br/>
<img src="./react/react-sample/src/images/MANAMANA_機能紹介.003.jpeg">
<br/>
<img src="./react/react-sample/src/images/MANAMANA_機能紹介.004.jpeg">
<br/>
<br/>
<h1>②どのような技術を用いて</h1>

◯バックエンド：Go 1.17

（選定理由）
実行速度が速いとの事で選定しました。<br/>
また、SPA形式でのアプリケーションの作成を考えていたので、<br/>
バックエンドのAPIとして使用するのに相性が良いと考え選定しました。
<br/>
<br/>
◯FW：echo v4

（選定理由）
Ginと悩みましたが、
Ginよりも軽量かつ小・中規模のアプリケーションの開発に向いているとの事で、<br/>
こちらを選定しました。<br/>
また、情報ソースも多く、<br/>
一見した時にコードの可読性が良かったのも選定理由の一つです。
<br/>
<br/>
◯ORM：Gorm

（選定理由）
保守性向上の為ORMの導入は元より考えていたので、<br/>
情報ソースが多く可読性の良いこちらを選定しました。
<br/>
<br/>
◯OSS：agouti

（選定理由）
当初ブラウザのクローリングはSeleniumで行う予定でしたが、<br/>
PythonやJavaの動作環境をその為だけに用意するのは保守性の面から鑑み、<br/>
あまり妥当ではないと考え、<br/>
Goで動作するこちらを選定しました。
<br/>
<br/>
◯フロントエンド：React 17.0.2

（選定理由）
Goと同様に実行速度が速いとの事で選定しました。<br/>
また、レンダリングも速く、<br/>
UXの面から鑑みてもこちらがフロントエンドに向いていると考え選定しました。
<br/>
<br/>
◯ライブラリ：axios

（選定理由）
SPAでアプリケーションを制作するにあたって非同期通信は必須となるので、<br/>
こちらを選定しました。<br/>
ajaxの使用経験はあったので、<br/>
学習コスト自体は低かったです。
<br/>
<br/>
◯Webサーバーソフトウェア：nginx 1.21.4

（選定理由）
Apacheと違い、<br/>
シングルスレッドでプロセスの処理を行う事ができ、<br/>
負荷の増大に寄与せず、<br/>
処理速度の維持が実現できるので、<br/>
こちらを選定しました。
<br/>
<br/>
◯インフラ管理ツール：Terraform v1.1.3

（選定理由）
インフラの保守性の向上から選定しました。<br/>
また、今後別のアプリケーションの開発を行う際に、<br/>
類似した環境構築が容易に再現できる様になるので、<br/>
導入するに至りました。
<br/>
<br/>
◯ツール：Dcoker 20.10.10

（選定理由）
開発環境構築の高速化、
およびローカル環境から本番環境への再現流用性の高さから選定しました。
<br/>
<br/>
<h1>③どのような工夫をして</h1>
【行動の背景】
<br/>
UXの向上を鑑みた時に、<br/>
SPA化が現状のベストプラクティスであると考え、<br/>
それを実現する為に技術を選定していきました。
<br/>
<br/>
<h1>④どのような成果に繋がったのか</h1>
これまでは、<br/>
自分の分かる範囲で使用経験のある技術を選定して開発を行なっていた事が多かったですが、<br/>
今回はどれも使用経験の無い技術で、<br/>
また、台頭間も無いモダンな技術もいくつか盛り込んでいるので、<br/>
そもそもWeb上に情報ソースが少なく、<br/>
開発中に問題が起きた際に解決するまでに時間が掛かってしまう事態が散見されました。
<br/>
<br/>
成果として、<br/>
まだ開発中の段階ではありますが、<br/>
自分の欲しい機能と使いたかった技術を盛り込み、<br/>
形にする事が出来てきたので、<br/>
更なる改良・改善を続けていきたいと思います。