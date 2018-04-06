package crawler

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/samjohnduke/crawl3/shared"

	"github.com/PuerkitoBio/goquery"
)

func TestExtractWithObject(t *testing.T) {
	for _, et := range extractTests {
		body := et.body
		var ej shared.Host
		err := json.Unmarshal([]byte(et.extracterJSON), &ej)
		if err != nil {
			t.Error(err)
		}

		doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(body)))
		if err != nil {
			t.Error(err)
		}

		extractorJSON := JSONExtractor{Rules: ej.Extractor}
		_, err = extractorJSON.Extract(doc)
		if err != nil {
			t.Error(err)
		}
	}
}

type extractTest struct {
	body          string
	extracterJSON string
}

var extractTests = []extractTest{
	extractTest{
		extracterJSON: `
		{
			"extractor": [{
				"@type": "NewsArticle",
				"@pageMatcher": [
					".news"
				],
				"fields": {
					"title": {
						"type": "String",
						"matcher": "h1",
						"content": "innerHTML"
					},
					"published_time": {
						"type": "Time",
						"matcher": "[itemprop=\"datePublished\"]",
						"content": "content"
					},
					"images": {
						"type": "[]String",
						"matcher": "#article-primary img",
						"content": "src"
					},
					"content": {
						"type": "String",
						"matcher": "#article-body p",
						"content": "innerHTML"
					}
				}
			}]
		}
	`,
		body: `
	
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <script type="text/javascript">window.NREUM||(NREUM={});NREUM.info = {"beacon":"bam.nr-data.net","errorBeacon":"bam.nr-data.net","licenseKey":"ae432a7c36","applicationID":"9841687","transactionName":"NFxRMURUWBBVBxJRDQ0WcjZmGlgGQxdJWRAXUFAJUxtXEEQc","queueTime":0,"applicationTime":1,"agent":"","atts":""}</script><script type="text/javascript">window.NREUM||(NREUM={}),__nr_require=function(e,t,n){function r(n){if(!t[n]){var o=t[n]={exports:{}};e[n][0].call(o.exports,function(t){var o=e[n][1][t];return r(o||t)},o,o.exports)}return t[n].exports}if("function"==typeof __nr_require)return __nr_require;for(var o=0;o<n.length;o++)r(n[o]);return r}({1:[function(e,t,n){function r(){}function o(e,t,n){return function(){return i(e,[f.now()].concat(u(arguments)),t?null:this,n),t?void 0:this}}var i=e("handle"),a=e(2),u=e(3),c=e("ee").get("tracer"),f=e("loader"),s=NREUM;"undefined"==typeof window.newrelic&&(newrelic=s);var p=["setPageViewName","setCustomAttribute","setErrorHandler","finished","addToTrace","inlineHit","addRelease"],d="api-",l=d+"ixn-";a(p,function(e,t){s[t]=o(d+t,!0,"api")}),s.addPageAction=o(d+"addPageAction",!0),s.setCurrentRouteName=o(d+"routeName",!0),t.exports=newrelic,s.interaction=function(){return(new r).get()};var m=r.prototype={createTracer:function(e,t){var n={},r=this,o="function"==typeof t;return i(l+"tracer",[f.now(),e,n],r),function(){if(c.emit((o?"":"no-")+"fn-start",[f.now(),r,o],n),o)try{return t.apply(this,arguments)}catch(e){throw c.emit("fn-err",[arguments,this,e],n),e}finally{c.emit("fn-end",[f.now()],n)}}}};a("setName,setAttribute,save,ignore,onEnd,getContext,end,get".split(","),function(e,t){m[t]=o(l+t)}),newrelic.noticeError=function(e){"string"==typeof e&&(e=new Error(e)),i("err",[e,f.now()])}},{}],2:[function(e,t,n){function r(e,t){var n=[],r="",i=0;for(r in e)o.call(e,r)&&(n[i]=t(r,e[r]),i+=1);return n}var o=Object.prototype.hasOwnProperty;t.exports=r},{}],3:[function(e,t,n){function r(e,t,n){t||(t=0),"undefined"==typeof n&&(n=e?e.length:0);for(var r=-1,o=n-t||0,i=Array(o<0?0:o);++r<o;)i[r]=e[t+r];return i}t.exports=r},{}],4:[function(e,t,n){t.exports={exists:"undefined"!=typeof window.performance&&window.performance.timing&&"undefined"!=typeof window.performance.timing.navigationStart}},{}],ee:[function(e,t,n){function r(){}function o(e){function t(e){return e&&e instanceof r?e:e?c(e,u,i):i()}function n(n,r,o,i){if(!d.aborted||i){e&&e(n,r,o);for(var a=t(o),u=m(n),c=u.length,f=0;f<c;f++)u[f].apply(a,r);var p=s[y[n]];return p&&p.push([b,n,r,a]),a}}function l(e,t){v[e]=m(e).concat(t)}function m(e){return v[e]||[]}function w(e){return p[e]=p[e]||o(n)}function g(e,t){f(e,function(e,n){t=t||"feature",y[n]=t,t in s||(s[t]=[])})}var v={},y={},b={on:l,emit:n,get:w,listeners:m,context:t,buffer:g,abort:a,aborted:!1};return b}function i(){return new r}function a(){(s.api||s.feature)&&(d.aborted=!0,s=d.backlog={})}var u="nr@context",c=e("gos"),f=e(2),s={},p={},d=t.exports=o();d.backlog=s},{}],gos:[function(e,t,n){function r(e,t,n){if(o.call(e,t))return e[t];var r=n();if(Object.defineProperty&&Object.keys)try{return Object.defineProperty(e,t,{value:r,writable:!0,enumerable:!1}),r}catch(i){}return e[t]=r,r}var o=Object.prototype.hasOwnProperty;t.exports=r},{}],handle:[function(e,t,n){function r(e,t,n,r){o.buffer([e],r),o.emit(e,t,n)}var o=e("ee").get("handle");t.exports=r,r.ee=o},{}],id:[function(e,t,n){function r(e){var t=typeof e;return!e||"object"!==t&&"function"!==t?-1:e===window?0:a(e,i,function(){return o++})}var o=1,i="nr@id",a=e("gos");t.exports=r},{}],loader:[function(e,t,n){function r(){if(!x++){var e=h.info=NREUM.info,t=d.getElementsByTagName("script")[0];if(setTimeout(s.abort,3e4),!(e&&e.licenseKey&&e.applicationID&&t))return s.abort();f(y,function(t,n){e[t]||(e[t]=n)}),c("mark",["onload",a()+h.offset],null,"api");var n=d.createElement("script");n.src="https://"+e.agent,t.parentNode.insertBefore(n,t)}}function o(){"complete"===d.readyState&&i()}function i(){c("mark",["domContent",a()+h.offset],null,"api")}function a(){return E.exists&&performance.now?Math.round(performance.now()):(u=Math.max((new Date).getTime(),u))-h.offset}var u=(new Date).getTime(),c=e("handle"),f=e(2),s=e("ee"),p=window,d=p.document,l="addEventListener",m="attachEvent",w=p.XMLHttpRequest,g=w&&w.prototype;NREUM.o={ST:setTimeout,SI:p.setImmediate,CT:clearTimeout,XHR:w,REQ:p.Request,EV:p.Event,PR:p.Promise,MO:p.MutationObserver};var v=""+location,y={beacon:"bam.nr-data.net",errorBeacon:"bam.nr-data.net",agent:"js-agent.newrelic.com/nr-1071.min.js"},b=w&&g&&g[l]&&!/CriOS/.test(navigator.userAgent),h=t.exports={offset:u,now:a,origin:v,features:{},xhrWrappable:b};e(1),d[l]?(d[l]("DOMContentLoaded",i,!1),p[l]("load",r,!1)):(d[m]("onreadystatechange",o),p[m]("onload",r)),c("mark",["firstbyte",u],null,"api");var x=0,E=e(4)},{}]},{},["loader"]);</script>
    <title>How Westpac dealt with its first phishing attack - Security - iTnews</title>
    <meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=no">    
	<meta name="theme-color" content="#051936">
    <meta name="Description" content="The genesis of industry collaboration.">
    <meta name="Keywords" content="acsc,interbank,jcsc,phishing,security,westpac">
    <meta name="news_keywords" content="acsc, interbank, jcsc, phishing, security, westpac">
    
    
    <link rel="alternate" type="application/rss+xml" title="iTnews RSS Feed" href="/rss/rss.ashx" />    
    <link href="https://fonts.googleapis.com/css?family=Martel:400,300,700,900" rel="stylesheet" type="text/css">
    <link href="https://fonts.googleapis.com/css?family=Lato:300,300italic,400,400italic,700" rel="stylesheet" type="text/css">
    
    
    <meta property="og:site_name" content="iTnews"> 
    <meta property="og:title" content="How Westpac dealt with its first phishing attack">
    <meta property="og:url" content="http://www.itnews.com.au/news/how-westpac-dealt-with-its-first-phishing-attack-487821">
    <meta property="og:image" content="https://i.nextmedia.com.au/News/20180327022945_phishing.jpg">
    <meta property="og:description" content="The genesis of industry collaboration.">            
    
    
    <meta name="twitter:card" content="summary">
    <meta name="twitter:site" content="@iTnews_au">
    <meta name="twitter:creator" content="@iTnews_au">
    <meta name="twitter:title" content="How Westpac dealt with its first phishing attack">
    <meta name="twitter:description" content="The genesis of industry collaboration.">
    <meta name="twitter:url" content="http://www.itnews.com.au/news/how-westpac-dealt-with-its-first-phishing-attack-487821">
    
        <meta name="twitter:image" content="https://i.nextmedia.com.au/News/20180327022945_phishing.jpg">
    
    
        
    
    <link rel="canonical" href="https://www.itnews.com.au/news/how-westpac-dealt-with-its-first-phishing-attack-487821"/>
    
    
      
    <link rel="stylesheet" href="https://ajax.googleapis.com/ajax/libs/jqueryui/1.11.4/themes/smoothness/jquery-ui.css">
    <link rel="stylesheet" type="text/css" href="/styles/css_258bb9db546fc997ead21e5f02deacb2.css" />  
        
    <!--[if lt IE 9]>
        <script src="//cdnjs.cloudflare.com/ajax/libs/html5shiv/3.6.2/html5shiv.js"></script>
        <script src="//s3.amazonaws.com/nwapi/nwmatcher/nwmatcher-1.2.5-min.js"></script>
        <script src="//cdn.jsdelivr.net/selectivizr/1.0.3b/selectivizr.js"></script>
        <script src="//cdnjs.cloudflare.com/ajax/libs/respond.js/1.1.0/respond.min.js"></script>
        <script src="/scripts/rem.min.js"></script>
    <![endif]-->
    
    
</head>
<body>    
<div class="wrapper news" id="wrapper">
    <div class="off-canvas-wrap" data-offcanvas>
    <div class="inner-wrap">
        <header>
            
                <div id="breaking-news" class="row full-width collapse show-for-large-up fixed">
                    <a id="breaking-news-header" style="float:left" href="/news">Latest News</a>
                    <div id="breaking-news-container">
                        
                                <a href="/news/nbn-co-opts-for-low-key-fttc-launch-488000" class="breaking-news-item">
                                    <h2 class="article-headline">NBN Co opts for low-key FTTC launch</h2>
                                </a>
                            
                                <a href="/news/internode-sets-nbn-peak-speed-benchmark-487971" class="breaking-news-item">
                                    <h2 class="article-headline">Internode sets NBN peak speed benchmark</h2>
                                </a>
                            
                                <a href="/news/uber-brings-carpool-service-to-australia-487938" class="breaking-news-item">
                                    <h2 class="article-headline">Uber brings carpool service to Australia</h2>
                                </a>
                            
                                <a href="/news/nbn-speeds-blitz-all-expectations-487939" class="breaking-news-item">
                                    <h2 class="article-headline">NBN speeds blitz all expectations</h2>
                                </a>
                            
                                <a href="/news/accc-puts-myobs-reckon-acquisition-in-doubt-487940" class="breaking-news-item">
                                    <h2 class="article-headline">ACCC puts MYOB's Reckon acquisition in doubt</h2>
                                </a>
                            
                            
                    </div>                                
                </div>
            
            <div class="collapse header-full show-for-large-up" id="header-secondary">
                <div class="header-columns">
                    <div id="header-lb">
                        <div id="ph-leader"></div>
                        <div id="div-gpt-ad-leader"></div>
                    </div> 
                    <a href="/" id="site-logo"></a> 
                </div>
            </div>
            <div class="collapse header-full" id="header-primary">                
                

<div class="header-columns">
<div class="sticky" id="sticky-nav">
    <nav class="top-bar row full-width collapse show-for-large-up" data-topbar role="navigation">
        <div id="top-bar-secondary">
            <div id="follow-us">
                <a class="follow-twitter" href="https://twitter.com/itnews_au" target="_blank"></a>
                <a class="follow-facebook" href="https://www.facebook.com/iTnewsAustralia" target="_blank"></a>
                <a class="follow-linkedin" href="https://www.linkedin.com/company/itnews" target="_blank"></a>
                <a class="follow-rss" href="/rss"></a>
            </div>
                
                <a class="login-click">LOG IN</a>
                <a class="subscribe" href="/register">SUBSCRIBE</a>
            
            
             <a href="#" class="search">&nbsp;</a>
        </div>
        <section class="top-bar-section">
            <a href="/" id="itnews-logo-sticky"><img src="/images/itnews-logo-sticky.png" alt="iTnews" /></a>
            <ul class="top-bar-primary">
                
                    <li><a href="/tag/governmentit">GOVERNMENT IT</a></li>
                
                    <li><a href="/technology/security">SECURITY</a></li>
                
                    <li><a href="/tag/financeit">FINANCE IT</a></li>
                
                    <li><a href="/technology/telco-isp">TELCO</a></li>
                
                    <li><a href="https://www.itnews.com.au/awards">BENCHMARK AWARDS</a></li>
                
            </ul> 
        </section>               
        <div id="search-box">
            <input type="text" placeholder="Search iTnews"><button>Search</button>
        </div>        
    </nav>
    
    
    <div class="hide-for-large-up">
    <nav class="mobile-nav" role="navigation">
     
        <section class="title-area full-width">
            <div class="menu-primary" >
                <button class="c-hamburger c-hamburger--htx">
                    <span></span>
                </button>
            </div>
            <div class="logo"><a href="/"><img src="/images/itnews-logo-white.png" alt="iTnews" /></a></div>
            <div class="menu-secondary" >
                <a class="user-icon"><div id="icon"></div></a>
            </div>
        </section>

        <section class="top-bar-nav">
            <div class="mobile-nav-search">
                <input type="text" placeholder="Search iTnews">
            </div>
            <div class="mobile-nav-container">
                
                <a href="/tag/governmentit" class="mobile-nav-item"><span>GOVERNMENT IT</span></a>
                
                <a href="/technology/security" class="mobile-nav-item"><span>SECURITY</span></a>
                
                <a href="/tag/financeit" class="mobile-nav-item"><span>FINANCE IT</span></a>
                
                <a href="/technology/telco-isp" class="mobile-nav-item"><span>TELCO</span></a>
                
                <a href="https://www.itnews.com.au/awards" class="mobile-nav-item"><span>BENCHMARK AWARDS</span></a>
                
            </div>
            <div class="mobile-nav-follow">
                <a class="follow-twitter" href="https://twitter.com/itnews_au" target="_blank"></a>
                <a class="follow-facebook" href="https://www.facebook.com/iTnewsAustralia" target="_blank"></a>
                <a class="follow-linkedin" href="https://www.linkedin.com/company/itnews" target="_blank"></a>
                <a class="follow-rss" href="/rss"></a>
            </div>
        </section>
        
        <section class="top-bar-user">
            
                <div id="login-form-mobile">                    
                    <h3 class="section-header"><span>Log In</span></h3>   
                    <div class="mobile-small">           
                        <div class="form-label">Username:</div>
                        <div class="form-input"><input id="username-mobile" name="username" type="text" required /></div>
                    </div>
                    <div class="mobile-small"> 
                        <div class="form-label">Password:</div>
                        <div class="form-input"><input id="password-mobile" name="password" type="password" required /></div>
                    </div>
                    <div class="row collapse form-checkbox">
                        <input id="rememberMe-mobile" name="rememberMe" type="checkbox" /><label for="rememberMe">Remember me</label>
                        <span>| &nbsp;<a href="/forgot" title="Forgot your password?">Forgot password?</a></span>
                    </div>
                    <div class="ui-dialog-buttonset mobile-submit"><input type="button" id="login-button-mobile" value="Log In" class="button" /></div>
                    <div id="login-validation-mobile"></div>
                    <div id="login-response-mobile"></div>                
                    <div id="mobile-register"><a href="/register">Don't have an account? Register now!</a></div>                
                </div>
            
            
        </section>
    </nav>
    </div>
    
    
</div>
</div>
            </div>
        </header>
        <div id="skin">
            <div class="row leaderboard">
                <div id="div-gpt-ad-pushdown"></div>
            </div>
            <div class="row main-section collapse">
                <div class="large-12 columns">
                    

<div class="row">
    

<article class="small-12 large-12 columns" itemscope itemtype="http://schema.org/NewsArticle">

    <h1 id="article-headline" itemprop="headline">How Westpac dealt with its first phishing attack</h1>
    <div id="article-details-mobile">
        By 
                <a href="/author/allie-coyne-461593">
                    <span itemprop="author" itemscope itemtype="http://schema.org/Person">
                        <span itemprop="name">Allie Coyne</span>
                    </span>
                </a>
                
         on <span id="ctl00_ContentPlaceHolder1_ucArticle_lbDateFromMobile" itemprop="datePublished" content="2018-03-28T06:00">Mar 28, 2018 6:00AM</span>
    </div>
    <div class="row collapse">
        <div id="article-primary">
            <div id="article-image" itemprop="image" itemscope itemtype="https://schema.org/ImageObject">
                <img id="ctl00_ContentPlaceHolder1_ucArticle_imgImage" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=https%3a%2f%2fi.nextmedia.com.au%2fNews%2f20180327022945_phishing.jpg&amp;w=480&amp;c=0&amp;s=1" alt="How Westpac dealt with its first phishing attack" style="border-width:0px;" />
                    <meta itemprop="url" content="https://i.nextmedia.com.au/News/20180327022945_phishing.jpg">
                    <meta itemprop="width" content="960">
                    <meta itemprop="height" content="540">
                    
                    
                
                
            
                
                <div id="linked-gallery-placeholder"></div>  
            </div>
            
            <div id="div-gpt-ad-island" class="mrec-top"></div>
            
            
            <h2 id="article-intro" itemprop="description">The genesis of industry collaboration.
</h2>
            
                       
                
            <div id="article-body" itemprop="articleBody">     
                <p><span style="font-weight: 400;">On Friday July 4 2003, Westpac chief information security officer Richard Johnson got a call from a concerned executive. Someone had created a fake Westpac online banking website.</span></p>

<p><span style="font-weight: 400;">It was the first time the bank had ever come face to face with what we&rsquo;d now consider an online banking phishing attack.</span></p>

<p><span style="font-weight: 400;">In 2003, the main kinds of security threats banks had to deal with were website defacements or mass mailer worms like VBS.</span></p>

<p><span style="font-weight: 400;">This fake banking site was new - and that made it highly effective.</span></p>

<p><span style="font-weight: 400;">&ldquo;It fooled our customers, who like us at the time, didn&rsquo;t even have a name for this new phenomenon,&rdquo; Johnson told <em>iTnews</em> on the sidelines of the recent FST Media Future of Security summit.</span></p>

<p><span style="font-weight: 400;">The phishing site - or what Australia&rsquo;s banks then termed a &lsquo;ghost&rsquo; website - was not well constructed: the grammar was poor, the URL was clearly wrong, there was no SSL padlock, and it contained Cyrillic Russian characters.</span></p>

<p><span style="font-weight: 400;">&ldquo;Overall it was a pretty poor imitation. But it fooled people because it was different; it was really the first time we&rsquo;d seen true digital-based attacks for money that were trying to fool customers,&rdquo; Johnson said.</span></p>

<p><span style="font-weight: 400;">The fake Westpac site was followed by very similar phishing attacks on Australia&rsquo;s other banks within a matter of weeks.</span></p>

<p><span style="font-weight: 400;">&ldquo;We all realised very quickly that all our customers were being targeted by the same threats. That was new. The bad guys had changed the rules of the game,&rdquo; Johnson said.</span></p>

<p><span style="font-weight: 400;">&ldquo;And we quickly learned that as an ecosystem, we could far better defend our customers against these criminals.&rdquo;</span></p>

<p><span style="font-weight: 400;">The immediate priority for Westpac and the other banks was to understand what they were dealing with.</span></p>

<p><span style="font-weight: 400;">Through rapid-fire phone calls between cyber security teams the banks managed to trace the attackers to a compromised server sitting in a Florida data centre.</span></p>

<p><span style="font-weight: 400;">Contacting the provider hosting the website to get it taken down, however, turned out to be more difficult - the site had been deliberately put up on the Independence Day long weekend for this very reason.</span></p>

<p><span style="font-weight: 400;">After a bit of effort the bank managed to get through to local tech support and get the site taken offline.&nbsp;</span><span style="font-weight: 400;">The next step was making sure that while the attackers might be able to get into customer accounts with the stolen credentials, they couldn&rsquo;t then send the money offshore.</span></p>

<p><span style="font-weight: 400;">&ldquo;Very quickly we took additional countermeasures to make that really hard - we delayed it, we put in extra visibility and controls, fraud systems in the backend, and really importantly we figured out that we needed to work with the federal police,&rdquo; Johnson said.</span></p>

<p><span style="font-weight: 400;">The banks reached out to the man who was then running the high-tech crime centre for the AFP, Alastair MacGibbon - who is now the head of the Australian Cyber Security Centre, and the cyber security advisor to the Prime Minister.</span></p>

<p><span style="font-weight: 400;">&ldquo;We as a group of banks worked with Alastair to work out what we needed to do help them get involved in dealing with this,&rdquo; Johnson said.</span></p>

<p><span style="font-weight: 400;">&ldquo;It meant that while we might not always be able to stop a phishing site pop up, we may be able to detect the criminals trying to move the money from one account to another to get it offshore.&rdquo;</span></p>

<p><strong>Interbank</strong></p>

<p><span style="font-weight: 400;">This set of events was the genesis for one of the most critical collaborations in the country today: the Interbank financial services industry cyber security network.</span></p>

<p><span style="font-weight: 400;">What initially started as frenzied phone calls in a time of crisis has morphed into a real-time online sharing network and quarterly meeting that has been in place for more than a decade.</span></p>

<p>The quarterly two-day event spans security, risk, and fraud and includes sessions from external speakers like cyber security experts in the federal government.</p>

<p><span style="font-weight: 400;">The events are backed by an online network that is used to share information back and forth at &ldquo;machine speed&rdquo;, according to Johnson.</span></p>

<p><span style="font-weight: 400;">&ldquo;It&rsquo;s a really great example of a collaborative network of white hats working together against that common enemy,&rdquo; he said.</span></p>

<p><span style="font-weight: 400;">&ldquo;We have learned that security is not a competitive advantage. Sharing intel and lessons learned can only strengthen the overall ecosystem and help make it harder on the attackers.&rdquo;</span></p>

<p><strong>Cyber security hubs</strong></p>

<p><span style="font-weight: 400;">It&rsquo;s this approach Johnson and Westpac are hoping will be replicated in the federal government-backed joint cyber security centres currently&nbsp;<a href="https://www.itnews.com.au/news/govt-opens-sydney-cyber-threat-sharing-centre-487401" rel="noopener" target="_blank">springing up around the country.</a></span></p>

<p><span style="font-weight: 400;">The cyber security industry in Australia is heavily concentrated around banks, telcos, and the federal government; the joint centres are intended to lower the barrier to entry for other organisations with less cyber security muscle.</span></p>

<p><span style="font-weight: 400;">They act like regional &ldquo;hubs&rsquo; for the Australian Cyber Security Centre and allow organisations across any industry to access and disseminate information across the network. </span></p>

<p><span style="font-weight: 400;">The centres physically co-locate</span><span style="font-weight: 400;">&nbsp;government, business, and academic cyber security experts and are&nbsp;complemented by an information-sharing portal currently being built.</span></p>

<p><span style="font-weight: 400;">The initiative stemmed from complaints about the ability to access the experts housed within the ACSC in the ASIO building in Canberra.</span></p>

<p><span style="font-weight: 400;">&ldquo;Most industry-based cyber security people are not in Canberra. This creates a place that we can come to and work together,&rdquo; Johnson said.</span></p>

<p><span style="font-weight: 400;">&quot;When we look ahead to threats that could be systemic or affect us as a nation - things like WannaCry and Petya - I can see before long the JCSCs could become the mechanism that we use to create a national security operations centre, where everyone can dial in and co-ordinate what they&rsquo;re seeing.</span></p>

<p><span style="font-weight: 400;">&ldquo;The bit I really hope to see the JCSC deliver is the incident response piece but also the outreach for the rest of the industry that just doesn&rsquo;t have that scale.</span></p>

<p><span style="font-weight: 400;">&ldquo;Imagine there&rsquo;s a critical event: my analysts are working with Telstra&rsquo;s and CBA&rsquo;s and NAB&rsquo;s and guys from Canberra. We figure out what to do, then we can develop a quick guide saying &#39;there&rsquo;s a new attack out there, here&rsquo;s what you need to know&#39;. So people who don&rsquo;t have those large teams can leverage off that capability.&rdquo;</span></p>

<p><span style="font-weight: 400;">Johnson thinks the joint cyber security centres put Australia in &ldquo;the best position we&rsquo;ve ever been in&rdquo;.</span></p>

<p><span style="font-weight: 400;">He noted the upward trend in malicious cyber activity - like the growing volume of phishing, malicious spam, and critical vulnerabilities - and argues it needs to be approached as a common enemy.</span></p>

<p><span style="font-weight: 400;">&ldquo;The whole experience in banking and finance against cybercrime has taught us the incredible value you can bring to play when you behave and think as an ecosystem,&rdquo; Johnson said.</span></p>

<p><span style="font-weight: 400;">&ldquo;For a country like Australia it&rsquo;s the right approach, because all of our citizens are customers of various entities and organisations. So it doesn&rsquo;t really matter which organisation was able to help protect that person, they&rsquo;re then protected. </span></p>

<p><span style="font-weight: 400;">&quot;There is a real benefit from that collaboration.&nbsp;</span><span style="font-weight: 400;">Collaboration is the single most important thing in the field of cyber security today.&rdquo;</span></p>

                
                
            </div>
               
            
            
            <div class="newsletter-subscribe-container">
                <a id="ctl00_ContentPlaceHolder1_ucArticle_hlNewsletterSubscribe" href="/register?src=ITN_ArticlePromo_1_487821"><img src="/images/newsletter-promo-1.png" style="border-width:0px;" /></a>
            </div>
              
            
            <div id="article-source">
                <a href="http://www.itnews.com.au" target="_blank">Copyright  &#169; iTnews.com.au</a> . All rights reserved. 
            </div>
            
            <div id="div-gpt-ad-inread"></div>
            <div id="article-tag-container">
                <div id="article-tag-header">Tags:</div><a class="article-tag" href="/tag/acsc">acsc</a> <a class="article-tag" href="/tag/interbank">interbank</a> <a class="article-tag" href="/tag/jcsc">jcsc</a> <a class="article-tag" href="/tag/phishing">phishing</a> <a class="article-tag" href="/technology/security">security</a> <a class="article-tag" href="/tag/westpac">westpac</a>
            </div>
            
                            
                <div id="article-security">
                    <a id="ctl00_ContentPlaceHolder1_ucArticle_hlSecurity" href="/technology/security">In Partnership With</a>
                </div>                
                 
            
            
            
        </div>
        <aside id="article-secondary">
            <div id="article-details">
                By 
                        <a href="/author/allie-coyne-461593">
                            <span itemprop="author" itemscope itemtype="http://schema.org/Person">
                                <span itemprop="name">Allie Coyne</span>
                            </span>
                        </a>
                        
                <br><span id="ctl00_ContentPlaceHolder1_ucArticle_lbDateFrom" itemprop="datePublished" content="2018-03-28T06:00">Mar 28 2018<br />6:00AM</span>
            </div>
            <div id="article-scroll">
                <div class="article-share">
                    <a class="share-comments" title="Join the Conversation" href="/news/how-westpac-dealt-with-its-first-phishing-attack-487821#disqus_thread" data-disqus-identifier="487821">0 Comments</a>
                    <a class="share-twitter" title="Share on Twitter" href="https://twitter.com/intent/tweet?text=How Westpac dealt with its first phishing attack&amp;related=&amp;url=https%3a%2f%2fwww.itnews.com.au%2fnews%2fhow-westpac-dealt-with-its-first-phishing-attack-487821%3Futm_source=desktop%26utm_medium=twitter%26utm_campaign=share" target="_blank"><div id="article-twitter-share-count"></div></a>
                    <a class="share-facebook" title="Share on Facebook" href="https://www.facebook.com/sharer/sharer.php?u=https%3a%2f%2fwww.itnews.com.au%2fnews%2fhow-westpac-dealt-with-its-first-phishing-attack-487821%3Futm_source=desktop%26utm_medium=facebook%26utm_campaign=share" target="_blank"><div id="article-facebook-share-count"></div></a>
                    <a class="share-linkedin" title="Share on LinkedIn" href="https://www.linkedin.com/shareArticle?url=https%3a%2f%2fwww.itnews.com.au%2fnews%2fhow-westpac-dealt-with-its-first-phishing-attack-487821%3Futm_source=desktop%26utm_medium=linkedin%26utm_campaign=share" target="_blank"><div id="article-linkedin-share-count"></div></a>
                    <a class="share-feedback" title="Send Feedback" href="/feedback/?id=487821"></a>
                    <a class="share-email" title="Email A Friend" href="mailto:?subject=iTnews: How Westpac dealt with its first phishing attack&amp;body=Check out this great article I read on iTnews:%0D%0A%0D%0AHow Westpac dealt with its first phishing attack%0D%0A%0D%0Ahttps%3a%2f%2fwww.itnews.com.au%2fnews%2fhow-westpac-dealt-with-its-first-phishing-attack-487821%3Futm_source=desktop%26utm_medium=email%26utm_campaign=share"></a>
                    <a class="share-print" title="Print" href="/tools/print.aspx?ciid=487821" target="_blank"></a>
                </div>
                
        <div>
            <h3 class="section-header"><span>Related Articles</span></h3>
            <ul id="related-article-list">
    
                <li><a href="/news/govt-opens-sydney-cyber-threat-sharing-centre-487401">Govt opens Sydney cyber threat sharing centre</a></li>
    
                <li><a href="/news/when-an-it-manager-falls-victim-to-a-phish-487280">When an IT manager falls victim to a phish</a></li>
    
                <li><a href="/news/major-cryptocurrency-exchange-hit-by-phishers-486663">Major cryptocurrency exchange hit by phishers</a></li>
    
                <li><a href="/news/google-offers-stronger-security-for-targeted-users-475619">Google offers stronger security for targeted users</a></li>
    
            </ul>
        </div>
    

        



     
            </div>
        </aside>
    </div>
    
    <div class="article-share-mobile">
        <a class="share-twitter" title="Share on Twitter" href="https://twitter.com/intent/tweet?text=How Westpac dealt with its first phishing attack&amp;related=&amp;url=https%3a%2f%2fwww.itnews.com.au%2fnews%2fhow-westpac-dealt-with-its-first-phishing-attack-487821%3Futm_source%3Dmobile%26utm_medium%3Dtwitter%26utm_campaign%3Dshare" target="_blank"><img src="/Images/mobile-share-twitter.png" border="0" alt="Share on Twitter" /></a>
        <a class="share-facebook" title="Share on Facebook" href="https://www.facebook.com/sharer/sharer.php?u=https%3a%2f%2fwww.itnews.com.au%2fnews%2fhow-westpac-dealt-with-its-first-phishing-attack-487821%3Futm_source%3Dmobile%26utm_medium%3Dfacebook%26utm_campaign%3Dshare" target="_blank"><img src="/Images/mobile-share-facebook.png" border="0" alt="Share on Facebook" /></a>
        <a class="share-linkedin" title="Share on LinkedIn" href="https://www.linkedin.com/shareArticle?url=https%3a%2f%2fwww.itnews.com.au%2fnews%2fhow-westpac-dealt-with-its-first-phishing-attack-487821%3Futm_source%3Dmobile%26utm_medium%3Dlinkedin%26utm_campaign%3Dshare" target="_blank"><img src="/Images/mobile-share-linkedin.png" border="0" alt="Share on LinkedIn" /></a>
        <a class="share-whatsapp" title="Share on WhatsApp" href="whatsapp://send?text=How Westpac dealt with its first phishing attack%0D%0A%0D%0Ahttps%3a%2f%2fwww.itnews.com.au%2fnews%2fhow-westpac-dealt-with-its-first-phishing-attack-487821%3Futm_source%3Dmobile%26utm_medium%3Dwhatsapp%26utm_campaign%3Dshare"><img src="/Images/mobile-share-whatsapp.png" border="0" alt="Share on Whatsapp" /></a>
        <a class="share-email" title="Email A Friend" href="mailto:?subject=iTnews: How Westpac dealt with its first phishing attack&amp;body=Check out this great article I read on iTnews:%0D%0A%0D%0AHow Westpac dealt with its first phishing attack%0D%0A%0D%0Ahttps%3a%2f%2fwww.itnews.com.au%2fnews%2fhow-westpac-dealt-with-its-first-phishing-attack-487821%3Futm_source%3Dmobile%26utm_medium%3Demail%26utm_campaign%3Dshare"><img src="/Images/mobile-share-email.png" border="0" alt="Email A Friend" /></a>
    </div>
    
    <div itemprop="publisher" itemscope itemtype="https://schema.org/Organization">
        <div itemprop="logo" itemscope itemtype="https://schema.org/ImageObject">
            <meta itemprop="url" content="https://www.itnews.com.au/images/itn-logo-clean.png">
            <meta itemprop="width" content="255">
            <meta itemprop="height" content="79">
        </div>
        <meta itemprop="name" content="iT News">
    </div>
</article>




        <div class="row small-12 columns" id="article-mostread-container">
            <h3 class="section-header"><span>Most Read Articles</span></h3>
            <div class="row collapse">
    
                <div class="small-12 medium-3 columns article-mostread-article">
                    <a href="/news/cbas-cio-out-in-restructure-487690" >
                        <img alt="CBA's CIO out in restructure" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=https%3a%2f%2fi.nextmedia.com.au%2fNews%2fdavid_whiteing.JPG&amp;h=271&amp;w=480&amp;c=1&amp;s=1">
                        <h2>CBA's CIO out in restructure</h2>
                    </a>
                </div>
    
                <div class="small-12 medium-3 columns article-mostread-article">
                    <a href="/news/nsw-transports-425m-it-overhaul-hits-the-skids-487707" >
                        <img alt="NSW Transport's $425m IT overhaul hits the skids" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=https%3a%2f%2fi.nextmedia.com.au%2fNews%2fNSW_trains_transport_(2).jpg&amp;h=271&amp;w=480&amp;c=1&amp;s=1">
                        <h2>NSW Transport's $425m IT overhaul hits the skids</h2>
                    </a>
                </div>
    
                <div class="small-12 medium-3 columns article-mostread-article">
                    <a href="/news/nbn-co-opts-for-low-key-fttc-launch-488000" >
                        <img alt="NBN Co opts for low-key FTTC launch" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=https%3a%2f%2fi.nextmedia.com.au%2fNews%2fnbn-fttc-deployment-halfway-03-1043.jpg&amp;h=271&amp;w=480&amp;c=1&amp;s=1">
                        <h2>NBN Co opts for low-key FTTC launch</h2>
                    </a>
                </div>
    
                <div class="small-12 medium-3 columns article-mostread-article">
                    <a href="/news/how-westpac-dealt-with-its-first-phishing-attack-487821" >
                        <img alt="How Westpac dealt with its first phishing attack" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=https%3a%2f%2fi.nextmedia.com.au%2fNews%2f20180327022945_phishing.jpg&amp;h=271&amp;w=480&amp;c=1&amp;s=1">
                        <h2>How Westpac dealt with its first phishing attack</h2>
                    </a>
                </div>
    
            </div>
        </div>
            


<div class="row">
    <div class="large-8 columns" id="disqus-container">

    
        <div id="disqus-login">
            You must be a registered member of <em>iTnews</em> to post a comment.<br />
            <a class="login-click">Log In</a> | 
            <a href="/register">Register</a>
        </div>
    

    
        <div id="disqus_thread"></div> 
        <noscript>Please enable JavaScript to view the <a href="http://disqus.com/?ref_noscript">comments powered by Disqus.</a></noscript>
    

    </div>
    
    <div class="large-4 columns">
        <div id="div-gpt-ad-island2" class="mrec-lower"></div>
        <div id="div-gpt-ad-button" class="button-ad-lower"></div>
        
        
        <section id="whitepapers-container">
            <h3 class="section-header"><span><a href="/whitepapers">Whitepapers from our sponsors</a></span></h3>
    
                <a class="row collapse whitepapers-list" href="/resource/the-business-case-for-hyperconvergence-478730" >
                    <div class="small-1 large-2 columns whitepapers-list-img"><img alt="The business case for hyperconvergence" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=https%3a%2f%2fi.nextmedia.com.au%2fReports%2f20171129015337_Nutanix_hyperconvergence_report_iTnews_2017.png&amp;w=100&amp;c=1&amp;s=0" class="width-only-img"></div>
                    <div class="small-11 large-10 columns whitepapers-list-title">The business case for hyperconvergence</div>
                </a>
    
                <a class="row collapse whitepapers-list" href="/resource/what-every-cio-should-know-about-devops-container-guides-by-puppet-476036" >
                    <div class="small-1 large-2 columns whitepapers-list-img"><img alt="What Every CIO Should Know about DevOps &amp; Container Guides by Puppet" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=http%3a%2f%2fi.nextmedia.com.au%2fReports%2fpuppet-logo.png&amp;w=100&amp;c=1&amp;s=0" class="width-only-img"></div>
                    <div class="small-11 large-10 columns whitepapers-list-title">What Every CIO Should Know about DevOps & Container Guides by Puppet</div>
                </a>
    
                <a class="row collapse whitepapers-list" href="/resource/the-5g-business-potential-industry-digitalisation-and-the-untapped-opportunities-for-operators-473599" >
                    <div class="small-1 large-2 columns whitepapers-list-img"><img alt="The 5G Business Potential &amp;#8211; Industry digitalisation and the untapped opportunities for operators" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=http%3a%2f%2fi.nextmedia.com.au%2fReports%2f20170919082034_Ericsson-.jpg&amp;w=100&amp;c=1&amp;s=0" class="width-only-img"></div>
                    <div class="small-11 large-10 columns whitepapers-list-title">The 5G Business Potential &#8211; Industry digitalisation and the untapped opportunities for operators</div>
                </a>
    
                <a class="row collapse whitepapers-list" href="/resource/solving-it-complexity-472486" >
                    <div class="small-1 large-2 columns whitepapers-list-img"><img alt="Solving IT complexity" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=http%3a%2f%2fi.nextmedia.com.au%2fReports%2fIBM-Logo.png&amp;w=100&amp;c=1&amp;s=0" class="width-only-img"></div>
                    <div class="small-11 large-10 columns whitepapers-list-title">Solving IT complexity</div>
                </a>
    
                <a class="row collapse whitepapers-list" href="/resource/optimising-enterprise-data-centres-for-the-cloud-465870" >
                    <div class="small-1 large-2 columns whitepapers-list-img"><img alt="Optimising Enterprise Data Centres for the Cloud" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=http%3a%2f%2fi.nextmedia.com.au%2fReports%2fjuniper-logo.jpg&amp;w=100&amp;c=1&amp;s=0" class="width-only-img"></div>
                    <div class="small-11 large-10 columns whitepapers-list-title">Optimising Enterprise Data Centres for the Cloud</div>
                </a>
    
        </section>
         
        
        <section id="events-container">
            <h3 class="section-header"><span><a href="/events">Events</a></span></h3>
            <ul>
    
                <li><a href="/event/acsc-conference-2018-487307">ACSC Conference 2018</a></li>
    
                <li><a href="/event/gartner-it-infrastructure-operations-management-data-center-summit-480736">Gartner IT Infrastructure, Operations Management & Data Center Summit</a></li>
    
                <li><a href="/event/future-of-mining-sydney-2018-486944">Future of Mining Sydney 2018</a></li>
    
                <li><a href="/event/digital-edge-experience-485979">Digital Edge Experience</a></li>
    
                <li><a href="/event/3rd-annual-design-thinking-summit-2018-487433">3rd Annual Design Thinking Summit 2018</a></li>
    
            </ul>
        </section>
         
        <div id="wout"><iframe src="https://haymarket-whistleout.s3.amazonaws.com/iTnews_Ad.html" width="300" height="250"></iframe></div>
    </div>
</div>
</div>


                
                    <div id="div-gpt-ad-footer" class="row footerboard"></div>
                    
                    
<script>
    ord=Math.random()*10000000000000000;
    document.write('<script src="/scripts/sponsoredcontent.ashx?type=SponsoredLink&si=Blogs&pa=&sc=32&output=script&ros=True&ord=' + ord + '"></scr' + 'ipt>');
</script>





                </div>
            </div>
        </div>
        
<div id="network-bar">
    <div class="row">
        <div class="large-12 columns">
            <h3 class="section-header">Most popular tech stories</h3>
        </div>
        <div class="large-12 columns">
        <ul class="large-block-grid-5 network-bar-grid">
            <li class="network-bar-item">
            <div class="network-bar-brand">&nbsp;</div>
            
                    <div>
                    <a href="http://www.crn.com.au/News/487834,telstra-turns-on-5g-over-wifi-for-gold-coast-locals.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">                   
                        
                        <img alt="Telstra turns on 5G-over-wifi for Gold Coast locals" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=https%3a%2f%2fi.nextmedia.com.au%2fNews%2fCRN_14_5G_stock.jpg&amp;h=141&amp;w=208&amp;c=1&amp;s=1">
                        
                        <h2 class="article-headline">
                            Telstra turns on 5G-over-wifi for Gold Coast locals
                        </h2>
                    </a>
                    </div>
                
                    <div>
                    <a href="http://www.crn.com.au/News/487934,nbn-speeds-from-iinet-optus-telstra-and-tpg-improving-but-still-held-back-by-fttn-accc.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">                   
                        
                        <h2 class="article-headline">
                            NBN speeds from iiNet, Optus, Telstra and TPG improving but still held back by FTTN: ACCC
                        </h2>
                    </a>
                    </div>
                
                    <div>
                    <a href="http://www.crn.com.au/News/487665,adelaide-based-telstra-partner-wireless-communications-acquired.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">                   
                        
                        <h2 class="article-headline">
                            Adelaide-based Telstra partner Wireless Communications acquired
                        </h2>
                    </a>
                    </div>
                
                    <div>
                    <a href="http://www.crn.com.au/News/487872,saps-small-business-cloud-service-sap-anywhere-is-shutting-down-only-had-30-active-customers.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">                   
                        
                        <h2 class="article-headline">
                            SAP's small business cloud service SAP Anywhere is shutting down, only had 30 active customers
                        </h2>
                    </a>
                    </div>
                
                    <div>
                    <a href="http://www.crn.com.au/News/487850,apple-unveils-new-ipad-education-software-to-win-back-schools.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">                   
                        
                        <h2 class="article-headline">
                            Apple unveils new iPad, education software to win back schools
                        </h2>
                    </a>
                    </div>
                
            </li>
            <li class="network-bar-item">
            <div class="network-bar-brand">&nbsp;</div>
            
                <div>
                    <a href="http://www.bit.com.au/Guide/464961,how-to-recover-deleted-emails-in-gmail.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                            <img alt="How to recover deleted emails in Gmail" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=http%3a%2f%2fi.nextmedia.com.au%2fFeatures%2frecover_gmail_1.jpg&amp;h=141&amp;w=208&amp;c=1&amp;s=1">
                        
                        <h2 class="article-headline">
                            How to recover deleted emails in Gmail
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.bit.com.au/News/337426,list-of-dates-when-australia-post-retail-outlets-will-be-closed-for-easter.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            List of dates when Australia Post retail outlets will be closed for Easter
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.bit.com.au/ReviewGroup/459180,6-cloud-accounting-systems-for-australian-small-businesses-compared-myob-quickbooks-reckon-saasu-sage-and-xero.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            6 cloud accounting systems for Australian small businesses compared: MYOB, QuickBooks, Reckon, Saasu, Sage and Xero
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.bit.com.au/Review/344651,7-accounting-packages-for-australian-small-businesses-compared-including-myob-quickbooks-online-reckon-xero.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            7 accounting packages for Australian small businesses compared: including MYOB, QuickBooks Online, Reckon, Xero
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.bit.com.au/News/316806,how-long-will-a-ups-keep-your-computers-on-if-the-lights-go-out.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            How long will a UPS keep your computers on if the lights go out?
                        </h2>
                    </a>
                </div>
                
            </li>
            <li class="network-bar-item">
            <div class="network-bar-brand">&nbsp;</div>
            
                <div>
                    <a href="http://www.pcauthority.com.au/Gallery/279957,top-25-fantasy-games-of-all-time.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                            <img alt="Top 25 fantasy games of all time" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=http%3a%2f%2fi.nextmedia.com.au%2fGalleries%2f20111114101345_dragon-age.jpg&amp;h=141&amp;w=208&amp;c=1&amp;s=1">
                        
                        <h2 class="article-headline">
                            Top 25 fantasy games of all time
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.pcauthority.com.au/Gallery/271730,top-15-obscure-video-game-consoles-for-collectors.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            Top 15 obscure video game consoles for collectors
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.pcauthority.com.au/Feature/462603,25-secret-whatsapp-tricks-you-probably-didnt-know-about.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            25 secret WhatsApp tricks you (probably) didn't know about
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.pcauthority.com.au/News/418496,how-to-remove-a-device-from-netflix-when-someones-accessing-your-account.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            How to: Remove a device from Netflix when someone's accessing your account
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.pcauthority.com.au/Feature/487890,hands-on-preview-huawei-p20-pro.aspx?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            Hands-on Preview: Huawei P20 Pro
                        </h2>
                    </a>
                </div>
                
            </li>
            <li class="network-bar-item">
            <div class="network-bar-brand">&nbsp;</div>
            
                    <div>
                        <a href="http://www.pcpowerplay.com.au/gallery/every-battlefield-game-ranked-from-worst-to-best,473272?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                            
                                <img alt="Every Battlefield game ranked from worst to best" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=http%3a%2f%2fi.nextmedia.com.au%2fGalleries%2f20170914095815_Battlefield-games-ranked-PC-PowerPlay011.jpg&amp;h=141&amp;w=208&amp;c=1&amp;s=1">
                            
                            <h2 class="article-headline">
                            Every Battlefield game ranked from worst to best
                        </h2>
                        </a>
                    </div>
                
                    <div>
                        <a href="http://www.pcpowerplay.com.au/feature/18-pro-tips-from-the-rainbow-six-siege-world-cup,450627?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                            
                            <h2 class="article-headline">
                            18 pro tips from the Rainbow Six Siege world cup
                        </h2>
                        </a>
                    </div>
                
                    <div>
                        <a href="http://www.pcpowerplay.com.au/feature/20-key-tips-for-succeeding-at-rainbow-six-siege,413238?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                            
                            <h2 class="article-headline">
                            20 key tips for succeeding at Rainbow Six Siege
                        </h2>
                        </a>
                    </div>
                
                    <div>
                        <a href="http://www.pcpowerplay.com.au/feature/far-cry-arcades-insane-potential-is-in-your-hands,487627?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                            
                            <h2 class="article-headline">
                            Far Cry Arcade&#8217;s insane potential is in your hands
                        </h2>
                        </a>
                    </div>
                
                    <div>
                        <a href="http://www.pcpowerplay.com.au/feature/10-advanced-tips-for-rainbow-six-siege,417453?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                            
                            <h2 class="article-headline">
                            10 advanced tips for Rainbow Six Siege
                        </h2>
                        </a>
                    </div>
                
            </li>
            <li class="network-bar-item">
            <div class="network-bar-brand">&nbsp;</div>
            
                <div>
                    <a href="http://www.iothub.com.au/news/startup-can-secure-iot-devices-with-data-analytics-487868?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                            <img alt="Startup can secure IoT devices with data analytics" src="https://i.nextmedia.com.au/Utils/ImageResizer.ashx?n=https%3a%2f%2fi.nextmedia.com.au%2fNews%2fiStock-696284840.jpg&amp;h=141&amp;w=208&amp;c=1&amp;s=1">
                        
                        <h2 class="article-headline">
                            Startup can secure IoT devices with data analytics
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.iothub.com.au/news/australian-govt-told-to-take-control-of-smart-cities-487748?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            Australian govt told to take control of smart cities
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.iothub.com.au/news/victorian-govt-to-pump-another-15m-into-agtech-487746?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            Victorian govt to pump another $15m into agtech
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.iothub.com.au/news/aws-to-talk-iot-at-sydney-summit-488043?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            AWS to talk IoT at Sydney summit
                        </h2>
                    </a>
                </div>
                
                <div>
                    <a href="http://www.iothub.com.au/news/edge-computing-will-be-key-to-iot-say-experts-488021?utm_source=itnews&amp;utm_medium=web&amp;utm_campaign=networkbar" target="_blank">
                        
                        <h2 class="article-headline">
                            Edge computing will be key to IoT, say experts
                        </h2>
                    </a>
                </div>
                
            </li>
        </ul>
        </div>
    </div>
</div>
        
<footer>
    <div id="footer-top-level" class="row">
        <div class="large-12 columns">
            <div id="footer-links">
                <a href="/page/contact">Contact Us</a>
                <a href="/page/about">About Us</a>
                
                <a href="/feedback">Feedback</a>
                <a href="/page/advertising">Advertise</a>
                <a href="/newsletter">Newsletter Archive</a>
                <a href="/sitemap">Site Map</a>
                <a href="/rss">RSS</a> 
            </div>
            <a href="http://www.nextmedia.com.au/" target="_blank" class="nm-logo"><img src="/images/logo_nextmedia.png" alt="nextmedia" /></a>&nbsp; <a href="/page/copyright">&copy; 2018 nextmedia Pty Ltd</a>.    
        </div>   
    </div>
    <div id="footer-bottom-level">
        <div class="row">
        <div class="small-12 medium-9 columns">
            <div id="footer-sites">
            <span>OTHER TECH SITES:</span>
            <a href="http://www.bit.com.au" target="_blank">BIT</a> &nbsp;|&nbsp;
            <a href="http://www.crn.com.au" target="_blank">CRN Australia</a> &nbsp;|&nbsp;
            <a href="http://www.iothub.com.au" target="_blank">IoT Hub</a> &nbsp;|&nbsp;
            <a href="http://www.pcauthority.com.au" target="_blank">PC &amp; Tech Authority</a> &nbsp;|&nbsp;
            <a href="http://www.pcpowerplay.com.au" target="_blank">PC PowerPlay</a>
            </div>
            <div id="footer-terms">
            All rights reserved. This material may not be published, broadcast, rewritten or redistributed in any form without prior authorisation.<br />Your use of this website
            constitutes acceptance of nextmedia's <a href="http://www.nextmedia.com.au/next-media-privacy-statement.html" target="_blank">Privacy Policy</a> and 
            <a href="http://www.nextmedia.com.au/next-media-terms-and-conditions.html" target="_blank">Terms & Conditions</a>.
            </div>
        </div>
        <div class="small-12 medium-3 columns" id="footer-rackspace"><a href="http://www.rackspace.com.au" target="_blank"><img src="/images/rackspace.png" alt="Powered by Rackspace Hosting" /></a></div>
        </div>
        
        <div id="otp">
            <div id="countdown" class="close"></div>
            <div id="div-gpt-ad-interstitial"></div>
        </div>
        <div class="mask"> </div>
        
        <div id="div-gpt-ad-skin"></div>
    </div>
</footer>
        

<div id="login-form">
    <form id="frm-login" action="/news/how-westpac-dealt-with-its-first-phishing-attack-487821" method="post">
        <h3 class="section-header"><span>Log In</span></h3>
        <div id="login-form-register"><a href="/register">Don't have an account? Register now!</a></div>
        <div id="login-validation"></div>
        <div id="login-response"></div>
        <div class="form-label">Username / Email:</div>
        <div class="form-input"><input id="username" name="username" type="text" required /></div>
        <div class="form-label">Password:</div>
        <div class="form-input"><input id="password" name="password" type="password" required /></div>
        <div class="row form-checkbox">
            <input id="rememberMe" name="rememberMe" type="checkbox" /><label for="rememberMe">Remember me</label><span>&nbsp; | &nbsp;<a href="/forgot" title="Forgot your password?">Forgot your password?</a></span>
        </div>
    </form>
</div>

    </div>
    <div id="itn-data-store"
        data-shareurl="%2fnews%2fhow-westpac-dealt-with-its-first-phishing-attack-487821"
        data-articleid="487821"
        data-adcategories="[&quot;security&quot;,&quot;security&quot;,&quot;technology&quot;]"
        data-adkeywords="[&quot;acsc&quot;,&quot;interbank&quot;,&quot;jcsc&quot;,&quot;phishing&quot;,&quot;security&quot;,&quot;westpac&quot;]"
        data-adsection="news"
        data-adshow="True"
        data-gasitesection="News"
        data-gasitecategory="Technology"
        data-gasitesubcategory="Security"
        data-gakeywords="|acsc|interbank|jcsc|phishing|security|westpac|"
        data-disqus_auth="bnVsbA== A30EC22D9FD9A6685CC3A85FBFE337C803A03AC6 1522662318"
        data-disqus_key="irryJPCUFvXWcfQC1pUkJNJgrmFUKbPzjmH2mIGajUKmrdpxmhw5yCENHGg0dbrf"
        data-disqus_shortname="itnewsnext"
        data-disqus_developer="0"
        data-disqus_title="How Westpac dealt with its first phishing attack"
        data-disqus_url="http://www.itnews.com.au/news/how-westpac-dealt-with-its-first-phishing-attack-487821"
        data-resizer="https://i.nextmedia.com.au/Utils/ImageResizer.ashx">
    </div>    
    
    <img id="imgPixelTrack" src="/t.ashx?u=&amp;c=487821&amp;s=3&amp;r=https%3a%2f%2fwww.itnews.com.au%2f&amp;n=%2fnews%2fArticle.aspx&amp;q=id%3d487821" alt="" style="display:none;" />
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.2/jquery.min.js"></script>
    <script src="https://ajax.googleapis.com/ajax/libs/jqueryui/1.11.4/jquery-ui.min.js"></script>
        
    <script src="/scripts/js_19dc54a312c7a7916a33906bc86ff67b.js"></script>    
        
    
    
    <!--[if lt IE 9]>    
        <link href="/styles/ie8.css" rel="stylesheet" type="text/css">
    <![endif]-->
    
    </div>   
</div>

</body>
</html>
`},
	extractTest{
		extracterJSON: `
    	{"extractor": [{
        	"@type": "NewsArticle",
            "@pageMatcher": [
        		".news.story_page"
        	],
            "fields": {
        		"title": {
        			"type": "String",
        			"matcher": "h1",
        			"content": "innerHTML"
        		},
        		"images": {
        			"type": "[]String",
        			"matcher": ".article img",
        			"content": "src"
        		},
        		"content": {
        			"type": "String",
        			"matcher": ".article",
        			"content": "innerHTML",
        			"exlcudeMatch": [
        				".tools",
        				".attached-content",
        				".byline",
        				".btn-group",
        				"h1",
        				".published",
        				".inline-content",
        				".topics",
        				".authorpromo",
        				".statepromo"
             		]
        		}
        	}
        }]}`,
		body: `
    <!DOCTYPE html>
  <!--[if lt IE 7]> <html lang="en" class="lt-ie9 lt-ie8 lt-ie7"> <![endif]-->
  <!--[if IE 7]> <html lang="en" class="lt-ie9 lt-ie8"> <![endif]-->
  <!--[if IE 8]> <html lang="en" class="lt-ie9"> <![endif]-->
  <!--[if gt IE 8]><!-->
  <html lang="en"> <!--<![endif]-->
  <head>
    <title>Stawell Gift: A case of sixth-time lucky for Victorian-based Tasmanian Jacob Despard - ABC News (Australian Broadcasting Corporation)</title>
  <meta http-equiv="Content-Type" content="text/html;charset=UTF-8"/>
<link rel="schema.DC" href="http://purl.org/dc/elements/1.1/"/>
<link rel="schema.DCTERMS" href="http://purl.org/dc/terms/"/>
<link rel="schema.iptc" href="urn:newsml:iptc.org:20031010:topicset.iptc-genre:8"/>
<link rel="amphtml" href="http://amp.abc.net.au/article/9610664"/>
<link rel="canonical" data-abc-platform="standard" href="http://www.abc.net.au/news/2018-04-02/tasmanian-jacob-despard-wins-the-stawell-gift-final/9610664"/>
                <meta name="DCSext.LocalRegion" content="westernvic"/><link rel="alternate" data-abc-platform="mobile" media="only screen and (max-width: 640px)" href="http://mobile.abc.net.au/news/2018-04-02/tasmanian-jacob-despard-wins-the-stawell-gift-final/9610664"/>
  <link media="all" rel="stylesheet" type="text/css"
                                                                              href="https://res.abc.net.au/bundles/2.2.0/styles/abc.bundle.2.2.0.min.css"/>
    <link rel="stylesheet" type="text/css"
                                                                              href="http://www.abc.net.au/res/sites/news-projects/news-core/1.3.3/desktop.css"/>
    <link rel="alternate" type="application/rss+xml" title="Just In" href="/news/feed/51120/rss.xml">
  <link rel="apple-touch-icon-precomposed" sizes="76x76" href="http://www.abc.net.au/cm/cb/8413652/News+iOS+76x76+2017/data.png"/>
    <link rel="apple-touch-icon-precomposed" sizes="120x120" href="http://www.abc.net.au/cm/cb/4355924/News+iOS+120x120/data.png"/>
    <link rel="apple-touch-icon-precomposed" sizes="152x152" href="http://www.abc.net.au/cm/cb/8413660/News+iOS+152x152+2017/data.png"/>
    <link rel="apple-touch-icon-precomposed" sizes="180x180" href="http://www.abc.net.au/cm/cb/8413674/News+iOS+180x180+2017/data.png"/>
    <script type="text/javascript" src="/news/ajax/5555684/managed.js"></script>
<script type="text/javascript" src="http://www.abc.net.au/res/libraries/jquery/jquery-1.11.3.min.js"></script>
        <script type="text/javascript" src="http://www.abc.net.au/res/libraries/location/abc.location-1.latest.min.js"></script>
        <script type="text/javascript" src="http://www.abc.net.au/res/bundles/platforms/abc.bundle.platforms-1.0.min.js"></script>
        <script type="text/javascript" src="http://www.abc.net.au/res/sites/news-projects/news-core/1.3.2/desktop.js"></script>
        <script type="text/javascript" src="/cm/code/8724582/abc.news.config-2018-03-20.js" async></script>
        <script type="text/javascript" src="/cm/code/5452002/Override+WebTrends+profile+for+Local+Radio+infosources+5.js"></script>
        </head>
  <body
          class="platform-standard news story_page">
  <!-- Start ABC Bundle Header 2.2.0 (customised) -->
<!--noindex-->
<nav id="abcHeader" class="global" aria-label="ABC Network Navigation" data-resourcebase="https://res.abc.net.au/bundles/2.2.0/" data-scriptsbase="https://res.abc.net.au/bundles/2.2.0/scripts/" data-version="2.2.0">
 <a class="abcLink home" href="http://www.abc.net.au/" data-mobile="http://mobile.abc.net.au/"><img src="https://res.abc.net.au/bundles/2.2.0/images/logo-abc@2x.png" width="65" height="16" alt="" />ABC Home</a>
 <div class="sites">
  <a class="controller" href="javascript:;" aria-controls="abcNavSites"><img src="https://res.abc.net.au/bundles/2.2.0/images/icon-menu-grey@1x.gif" 
   data-src="images/icon-menu-grey@1x.gif" data-hover="images/icon-menu-blue@1x.gif" class="icon" alt="" /><span class='text'><span>Open</span> Sites <span>menu</span></span></a>
  <div id="abcNavSites" class="menu" role="menu" aria-expanded="false">
   <ul>
    <li class='odd'><a class="abcLink" role="menuitem" href="http://www.abc.net.au/" data-mobile="http://mobile.abc.net.au/">ABC Home</a></li>
    <li><a class="abcLink" role="menuitem" href="http://www.abc.net.au/news/" data-mobile="http://mobile.abc.net.au/news/">News</a></li>
    <li class='odd'><a class="abcLink" role="menuitem" href="http://www.abc.net.au/iview/" data-mobile="http://www.abc.net.au/iview/">iview</a></li>
    <li><a class="abcLink" role="menuitem" href="http://www.abc.net.au/tv/" data-mobile="http://www.abc.net.au/tv/">TV</a></li>
    <li class='odd'><a class="abcLink" role="menuitem" href="http://www.abc.net.au/radio/" data-mobile="http://www.abc.net.au/radio/">Radio</a></li>
    <li><a class="abcLink" role="menuitem" href="http://www.abc.net.au/children/" data-mobile="http://mobile.abc.net.au/children/">Kids</a></li>
    <li class='odd'><a class="abcLink" role="menuitem" href="https://shop.abc.net.au/" data-mobile="https://shop.abc.net.au/">Shop</a></li>
    <li><a class="abcLink more" role="menuitem" href="http://www.abc.net.au/more/" data-mobile="http://www.abc.net.au/more/">More</a></li>
   </ul>
  </div>
 </div>
 <div class="accounts">
  <!-- Accounts is currently injected due to different URLs --><span data-src="images/icon-user-grey@1x.png" data-hover="images/icon-user-blue@1x.png"></span>
 </div>
 <a class="abcLink search" href="http://search.abc.net.au/s/search.html?collection=abcall_meta&form=simple"
  data-mobile="http://search.abc.net.au/s/search.html?collection=abcall_meta&form=simple"><span>Search</span>
  <img src="https://res.abc.net.au/bundles/2.2.0/images/icon-search-grey@1x.png" data-hover="images/icon-search-blue@1x.png" class="icon" alt="" /></a>
</nav>
<!--endnoindex-->
<!-- End ABC Bundle Header 2.2.0 (customised) -->
<!--noindex-->
  <div class="page_margins">
    <div id="header" class="header">
    <div class="brand">
        <a href="/news/"><img class="print" src="/cm/lb/8212706/data/news-logo-2017---desktop-print-data.png" alt="ABC News" width="387" height="100" />
                <img class="noprint" src="/cm/lb/8212704/data/news-logo2017-data.png" alt="ABC News" width="242" height="80" />
                </a></div>
      <a href="/news/australia/" class="location-widget">Australia</a><a href="/news/weather/" class="weather-widget">Weather</a></div>
  <div id="nav" class="nav">           
  <ul id="primary-nav">
        <li id="n-news" class=""><a href="/news/">News Home</a></li><li id="n-justin" class=""><a href="/news/justin/">Just In</a></li><li id="n-politics" class=""><a href="/news/politics/">Politics</a></li><li id="n-world" class=""><a href="/news/world/">World</a></li><li id="n-business" class=""><a href="/news/business/">Business</a></li><li id="n-sport" class=""><a href="/news/sport/">Sport</a></li><li id="n-science" class=""><a href="/news/science/">Science</a></li><li id="n-health" class=""><a href="/news/health/">Health</a></li><li id="n-arts-culture" class=""><a href="/news/arts-culture/">Arts</a></li><li id="n-analysis-and-opinion" class=""><a href="/news/analysis-and-opinion/">Analysis</a></li><li id="n-factcheck" class=""><a href="/news/factcheck/">Fact Check</a></li><li id="n-more" class=""><a href="/news/more/">More</a></li></ul>
  </div>
<!-- A modules - start -->
  <!-- A modules - end -->
  <div class="page section">
    <div class="subcolumns">
    <div class="c75l">
      <!-- B modules - start -->
      <div class="article section">
<div class="tools">  
  <a class="button"
     href="http://www2b.abc.net.au/EAF/View/MailToQuery.aspx?http://www.abc.net.au/news/2018-04-02/tasmanian-jacob-despard-wins-the-stawell-gift-final/9610664"><span>Email</span></a>
</div><!--endnoindex-->
  <h1>Stawell Gift: A case of sixth-time lucky for Victorian-based Tasmanian Jacob Despard</h1><div class="byline">
      <div class="bylinepromo">
          <a href="http://www.abc.net.au/westernvic/">ABC Western Vic</a>
        </div>
      By <a href="/news/9090134" target="_self" title="">Bridget Rollason</a></div>
    <p class="published">
      Updated 
    <span class="timestamp">
      April 02, 2018 16:19:36
    </span>
    </p>
  
<div class="inline-content photo full">
    <a href="/news/2018-04-02/stawell-gift-mens-winner-receives-his-cheque/9610696"> 
  <img src="http://www.abc.net.au/news/image/9610656-3x2-700x467.jpg" alt="A young man holds an oversized winner&#039;s cheque at the Stawell Gift running race in Western Victoria." title="Stawell Gift Men&#039;s winner Jacob Despard  receives his cheque" width="700" height="467"/>
</a><a href="/news/2018-04-02/stawell-gift-mens-winner-receives-his-cheque/9610696" class="inline-caption"><strong>
        Photo:</strong>
       Jacob Despard holds his winner's cheque after winning this year's Stawell Gift men's final. <span class="source">(ABC Western Victoria: Bridget Rollason)</span>
      </a></div>
<div class="attached-content">
      <div class="inline-content map left">
      <div class="story-map"></div>
      <a class="inline-caption" href="http://maps.google.com/?q=-37.0578,142.7796(Stawell%203380)&amp;z=5">
        <strong>Map: </strong>
        Stawell 3380</a>
    </div>
  </div>
  <p>Victorian-based Tasmanian Jacob Despard has taken out the men's final of the Stawell Gift, on his sixth attempt in a row at the title. </p><p>The 21-year-old held off last year's winner Matthew Rizzo and second-placed Melbourne runner Hamish Adams, to win the 137th Gift with a time of 12.12 seconds.</p><p>Despard beat US Olympian hurdler Devon Allen in the semi-finals and qualified for the final with the fastest time.</p><div class="inline-content photo full">
    <a href="/news/2018-04-02/stawell-gift-mens-winner-crosses-the-line/9610692">
  <img src="http://www.abc.net.au/news/image/9610650-3x2-700x467.jpg" alt="Male runners cross the finish line in the Stawell Gift in Western Victoria" title="Stawell Gift Men&#039;s winner Jacob Despard crosses the line" width="700" height="467"/>
</a><a href="/news/2018-04-02/stawell-gift-mens-winner-crosses-the-line/9610692" class="inline-caption"><strong>
        Photo:</strong>
       Tasmania's Jacob Despard wins the men's Stawell Gift final, beating last year's winner Matt Rizzo. <span class="source">(ABC Western Victoria: Bridget Rollason)</span>
      </a></div><p> </p><p>He was swamped by family and friends, who travelled across the country to watch him race, as soon as he crossed the finish line.</p><p>"My mother and my Nan have flown over from Tasmania and my Dad is from Western Australia," Despard said.</p><p>"They are just so supportive of me and I'm just so happy they made the effort to come and watch me today.</p><p>"It's always been my dream to win the Stawell Gift and I've tried for the last six years."</p><p>Despard said he planned to save most of his $40,000 winnings. </p><p>"$40,000 is a lot of money," he said.</p><p>"I've just moved to Melbourne and I spent a lot of money for the move.</p><blockquote class="quote--pullquote"><p>"I'm only 21, so I don't have the most money going around, so I'll save some to help me through the next six to 12 months, and then I think I'll spend a bit as well just to reward myself."</p></blockquote><p>Despard praised his coach Scott Rowsell for the win.</p><p>"I wouldn't like to know the amount of hours in the last six months to get me to where I am today — he's worked tirelessly," Despard said.</p><p>"The last six months have been going towards this race, I've been working so hard and to cross that line and to win, it's just unbelievable.</p><p>"It still hasn't sunk in yet I've got to say, it's just crazy."</p><p>Despard said he loved competing in the Stawell Gift and hoped to be back next year to defend his title.</p><h2>Female winner thanks her family</h2><p>Queensland university student Elizabeth Forsyth took out the women's Stawell Gift in a dramatic finish, falling over the finish line with a time of 13.69 seconds.</p><div class="inline-content photo full">
    <a href="/news/2018-04-02/womens-stawell-gift-winner-crosses-the-line/9610690">  
  <img src="http://www.abc.net.au/news/image/9610634-3x2-700x467.jpg" alt="Female runners cross the finish line of the Stawell Gift in Western Victoria." title="Women&#039;s Stawell Gift winner Elizabeth Forsyth crosses the line" width="700" height="467"/>
</a><a href="/news/2018-04-02/womens-stawell-gift-winner-crosses-the-line/9610690" class="inline-caption"><strong>
        Photo:</strong>
       Gold Coast surf lifesaver Elizabeth Forsyth wins the women's final of the Stawell Gift in Western Victoria. <span class="source">(ABC Western Victoria: Bridget Rollason)</span>
      </a></div><p> </p><p>The 22-year-old beat second-placed Pamela Austin with a time of 13.69 seconds.</p><p>The emotional athlete thanked her parents watching at home from the Gold Coast and her coach Brett Robinson.</p><p>"They've all rallied behind me this year and I can't thank them enough for all that they've done for me," Forsyth said.</p><p>The surf lifesaver also took home $40,000 prize money for her win.</p><div class="inline-content photo full">
    <a href="/news/2018-04-02/elizabeth-forsyth-holds-up-winners-cheque-for-stawell-gift-fin/9610698">  
  <img src="http://www.abc.net.au/news/image/9610678-3x2-700x467.jpg" alt="A young woman holds an oversized winner&#039;s cheque after winning the Stawell Gift women&#039;s final" title="Elizabeth Forsyth holds up winner&#039;s cheque for Stawell Gift final, April 2 2018" width="700" height="467"/>
</a><a href="/news/2018-04-02/elizabeth-forsyth-holds-up-winners-cheque-for-stawell-gift-fin/9610698" class="inline-caption"><strong>
        Photo:</strong>
       Gold Coast surf lifesaver Elizabeth Forsyth shows off her winner's cheque after the Stawell Gift women's final. <span class="source">(ABC Western Victoria: Bridget Rollason)</span>
      </a></div><p> </p>
<p class="topics">
	<strong>Topics:</strong>
	<a href="/news/topic/athletics">athletics</a>,	
	<a href="/news/topic/sport">sport</a>,
	<a href="/news/topic/stawell-3380">stawell-3380</a>
</p>
  <p class="published">
    First posted 
    <span class="timestamp">
      April 02, 2018 16:10:57
    </span>
  </p>
<div class="statepromo">
      <p>
<a href="/news/vic/" class="button">
<span><strong>
  More
</strong> stories from Victoria</span></a>

</p>
    </div>
  <!--noindex-->
  </div>
<!-- B modules - end -->
    </div>
    <div class="c25r sidebar">
    <!-- Sidebar modules - start -->
    <div class="section promo-list" style="background-image: url(/cm/lb/8784804/thumbnail/sidebar-thumbnail.png)">
        <h2  style="height: 60px">
          <a href="http://www.abc.net.au/westernvic">From ABC Western Vic</a></h2>

        <div class="inner">
         
<a href="http://www.abc.net.au/westernvic" class="button">
<span> More</span></a></div>
      </div>
    <div class="section localised-top-stories">
        <div class="inner">
          <h2>Top Stories</h2>
</div>
      </div>
    <div class="section graphic">
  <a href="http://m.me/abcnews.au?ref=welcomefromABCNewshomepage">
  <img src="/cm/lb/7996774/thumbnail.png" alt="Start a chat with ABC News on Facebook Messenger" title="Start a chat with ABC News on Facebook Messenger" class="" width="220" />
</a></div><div class="connect-with-abc section promo">
<h2>Connect with ABC News</h2>
<ul>
    <li><a href="http://www.facebook.com/abcnews.au" target="_self" title="">
        <span class="valign-helper"></span><img src="/cm/lb/6388890/data/facebook-data.png" alt="ABC News on Facebook" title="ABC News on Facebook" width="30" height="30" />
    </a></li>
    <li><a href="http://instagram.com/abcnews_au" target="_self" title="">
        <span class="valign-helper"></span><img src="/cm/lb/6389068/data/instagram-data.png" alt="ABC News on Instagram" title="ABC News on Instagram" width="32" height="31" />
    </a></li>
    <li><a href="http://www.twitter.com/abcnews" target="_self" title="">
        <span class="valign-helper"></span><img src="/cm/lb/6388894/data/twitter-data.png" alt="ABC News on Twitter" title="ABC News on Twitter" width="32" height="27" />
    </a></li>
    <li><a href="http://www.youtube.com/newsonabc" target="_self" title="">
        <span class="valign-helper"></span><img src="/cm/lb/6389052/data/youtube-data.png" alt="ABC News on YouTube" title="ABC News on YouTube" width="50" height="22" />
    </a></li>
    <li><a href="/news/2016-06-16/what-is-snapchat-and-why-is-the-abc-on-it/7517174" target="_self" title="">
        <span class="valign-helper"></span><img src="/cm/lb/8208710/data/snapchat-data.jpg" alt="ABC News on Snapchat" title="ABC News on Snapchat" width="31" height="30" />
    </a></li>

    <li><a href="https://apple.news/TsieN9U-WQ0SZbtP1VTI4uQ" target="_self" title="">
        <span class="valign-helper"></span><img src="/cm/lb/8458084/data/applenews-data.png" alt="ABC News on Apple News" title="ABC News on Apple News" width="32" height="32" />
    </a></li>
</ul>
<div class="clear"></div>
</div>
<div class="section graphic">
  <a href="/news/feeds/">
  <img src="/cm/lb/4745996/thumbnail/news-podcasts-sidebar-graphic-promo-thumbnail.jpg" alt="News Podcasts" title="News Podcasts" class="" width="220" />
</a></div><div class="section promo">
  <div class="inner">
      <h2>
</h2><p>If you have inside knowledge of a topic in the news, <a href="/news/contact/tip-off/" target="_self" title="">contact the ABC</a>.</p></div>
  </div>
<div class="related">
<div class="newsmail-signup card">
<h2>News in your inbox</h2>
<span class="spiel">Top headlines, analysis, breaking&nbsp;alerts</span>

<div class="form main">
</div>
<a href="/news/alerts/subscribe/">More info</a>
</div>
</div>
<div class="section promo">
  <div class="inner">
      <h2>
    <a href="http://www.abc.net.au/news/about/backstory/"
    >
        <span>ABC Backstory</span>
  </a>
</h2><p><a href="http://www.abc.net.au/news/about/backstory/" target="_self" title="">ABC teams share the story behind the story and insights into the making of digital, TV and radio content.</a></p></div>
  </div>
<div class="section promo">
  <div class="inner">
      <h2>

    <a href="http://about.abc.net.au/how-the-abc-is-run/what-guides-us/abc-editorial-standards/"
    >
          <span>Editorial Policies</span>
      </a>
</h2><p><a href="http://about.abc.net.au/how-the-abc-is-run/what-guides-us/abc-editorial-standards/" target="_self" title="">Read about our editorial guiding principles and the enforceable standard our journalists follow.</a></p></div>
  </div>
<!-- Sidebar modules - end -->
  </div>
</div>
</div><!-- C modules - start -->
  <div class="page section featured-scroller featured-scroller-4 dark">
<div class="section">
  <h2>Features</h2>
  </div>
<div class="inner">
  <ol class="subcolumns">
    <li class="c25l">
                  <div class="section">
		  <h2>
			 <a href="/news/2018-04-02/five-tips-to-help-you-make-the-most-of-reading-to-your-children/9609582" class="thumb">

  <img src="http://www.abc.net.au/news/image/8579586-16x9-220x124.jpg" alt="A man and male child sit on a couch reading a book together, you can&#039;t see their faces but the child is pointing at the pictures" title="Reading a book" width="220" height="124"/>
<span class="story"></span>
				<span class="label">
				  5 tips for reading to your children</span>
				<span class="border"></span>
			  </a></h2>
		 <p>Finding time to read to your children can be hard, but there are several ways you can make sure your child gets the most out of time for reading aloud.</p>
	 </div>
    </li>
         <li class="c25l">
                  <div class="section">
		  <h2>
			 <a href="/news/2018-04-01/australias-first-off-the-grid-solar-powered-home-40-years-on/9425648" class="thumb">
  <img src="http://www.abc.net.au/news/image/9425672-16x9-220x124.jpg" alt="The exterior of Michael and Judy Bos&#039;s house in Pearcedale." title="The exterior of Michael and Judy Bos&#039;s house" width="220" height="124"/>
<span class="story"></span>
				<span class="label">
				  Australia&#039;s first solar-powered home</span>
				<span class="border"></span>
			  </a></h2>
		 <p class="title">
        By <a href="/news/james-oaten/5512204" target="_self" title="">James Oaten</a></p>
		<p>Amid fears of an oil shortage in the 1970s, a Victorian couple decided they wanted to get off the electricity grid and build Australia's first fully solar-powered home.</p>
	 </div>
    </li>
         <li class="c25l">
                  <div class="section">
		  <h2>
			 <a href="/news/2018-04-02/paine-and-cummins-show-the-fight-australia-so-desperately-wants/9609562" class="thumb">
  <img src="http://www.abc.net.au/news/image/9609572-16x9-220x124.jpg" alt="Australia&#039;s captain Tim Paine plays a cut shot" title="Tim paine cuts" width="220" height="124"/>
<span class="story"></span>
				<span class="label">
				  Paine and Cummins dig in</span>
				<span class="border"></span>
			  </a></h2>
		 <p>After a top-order collapse, wicket-keeper Tim Paine and fast bowler Pat Cummins put on a stubborn 99-run stand in Johannesburg.</p>
	 </div>
    </li>
         <li class="c25r">
                  <div class="section">
		  <h2>
			 <a href="/news/2018-04-02/aboriginal-children-need-safe-culturally-appropriate-homes/9564006" class="thumb">

  <img src="http://www.abc.net.au/news/image/8260626-16x9-220x124.jpg" alt="Aboriginal woman and three little kids" title="Ms Walker and her grandchildren" width="220" height="124"/>
<span class="story"></span>
				<span class="label">
				  The homes Aboriginal children need</span>
				<span class="border"></span>
			  </a></h2>
		 <p>All children have the right to a safe, happy and emotionally supported childhood where they are nurtured and loved, but recent arguments about the removal of Indigenous children from their families fail to appreciate the complexity of the issue.</p>
	 </div>
    </li>
         </ol>
</div></div><div id="footer-stories" class="page section">
    <div class="subcolumns">
    <div class="c25l">
        <div class="section">
          <h2>Top Stories</h2>

</div>
      </div>
    <div class="c25l">
        <div class="section">
          <h2>Most Popular</h2>
          </div>
      </div>
    <div class="c25r">
        <div class="section">
          <h2>Analysis &amp; Opinion</h2>
</div>
      </div>
    </div>
  </div>
<!-- C modules - end -->
<div id="footer" class="page section">
    <!-- program footer-->
  <div class="subcolumns">
    <div id="sitemap" class="c75l">
        <div class="section">
          <h2>Site Map</h2>
        </div>
        <div class="subcolumns">
          <div class="c16l">
              <div class="section">
    <h3>Sections</h3>
    </div>
</div>
            <div class="c16l">
              <div class="section">
    <h3>Local Weather</h3>
    </div>
</div>
            <div class="c16l">
              <div class="section">
    <h3>Local News</h3>
    </div>
</div>
            </div>
      </div>

  </div>
</div>
  
  </body>
  </html>

    `,
	},
}
