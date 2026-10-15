# Dwaar

Dwaar (द्वार) हे SSH आणि लोकल पोर्ट फॉरवर्डिंग नियम व्यवस्थापित करण्यासाठी एक शक्तिशाली आणि सोपे Terminal User Interface (TUI) आहे. सुप्रसिद्ध [Bubble Tea](https://github.com/charmbracelet/bubbletea) फ्रेमवर्क वापरून Go मध्ये तयार केलेले, हे टर्मिनल ॲप्लिकेशन्सना Elm सारखी आर्किटेक्चर प्रदान करते, ज्यामुळे वापरकर्त्याचा अनुभव उत्तम आणि वेगवान होतो.

## विहंगावलोकन (Overview)

CLI द्वारे पोर्ट फॉरवर्डिंग नियम व्यवस्थापित करणे कंटाळवाणे होऊ शकते, विशेषतः एकाधिक कनेक्शन्स आणि SSH हॉप्ससह. Dwaar हे एक एकत्रित इंटरफेस प्रदान करून सोपे करते:
- लोकल पोर्ट ते रिमोट होस्ट:पोर्ट मॅपिंग कॉन्फिगर करणे.
- SSH टनेल्स सहजतेने सुरू करणे आणि थांबवणे.
- रिअल-टाइम कनेक्शन लॉग आणि टनेल स्थिती (status) पाहणे.
- पुढील वेळेस जलद सेटअपसाठी कॉन्फिगरेशन जतन करणे.

## वैशिष्ट्ये (Features)

- **TUI व्यवस्थापन**: अंगभूत फॉर्म इंटरफेस वापरून पोर्ट फॉरवर्डिंग नियम जोडा (Add), संपादित करा (Edit), आणि हटवा (Delete); SSH हॉप सेटिंग्ज आणि ऑटोस्टार्ट (autostart) टॉगलसह.
- **टनेल नियंत्रण**: एकाच कीस्ट्रोकने टनेल्स सुरू करा आणि थांबवा.
- **स्थिती निर्देशक (Status Indicators)**: `running`, `stopped`, `connecting`, आणि `error` स्टेट्ससाठी दृश्यमान संकेत.
- **लाइव्ह ट्रॅफिक आकडेवारी (Live Traffic Stats)**: टनेल चालू असताना प्रत्येक टनेलचे पाठवलेले/प्राप्त (sent/received) बाइट काउंटर्स.
- **कायमस्वरूपी कॉन्फिगरेशन**: नियम YAML कॉन्फिगरेशन फाइलमध्ये सेव्ह केले जातात आणि लोड केले जातात.
- **SSH सपोर्ट**: की आणि पासवर्ड प्रमाणीकरणासह (authentication) SSH वर पोर्ट फॉरवर्ड करा.
- **थीम्स (Themes)**: डीफॉल्ट आणि Catppuccin थीममध्ये रनटाइमवर बदल (switch) करा.

## तंत्रज्ञान (Tech Stack)

- **कोर फ्रेमवर्क**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **स्टाईलिंग**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **कंपोनंट्स**: [Bubbles](https://github.com/charmbracelet/bubbles)
- **फॉर्म्स**: [Huh](https://github.com/charmbracelet/huh)
- **SSH इंजिन**: `golang.org/x/crypto/ssh`

## इन्स्टॉलेशन (Installation)

तुमच्याकडे [Go](https://golang.org/doc/install) इन्स्टॉल असल्याची खात्री करा (आवृत्ती 1.26+ शिफारसीय).

```bash
git clone <repository-url>/dwaar
cd dwaar
go build -o dwaar ./cmd/dwaar
```

## वापर (Usage)

ॲप्लिकेशन सुरू करण्यासाठी बिल्ट एक्झिक्युटेबल रन करा:

```bash
./dwaar
```

### की-बाइंडिंग्ज (Key Bindings)

| की (Key) | कृती (Action) |
|---|---|
| `↑` / `↓` किंवा `j` / `k` | यादीत वर-खाली जाण्यासाठी |
| `/` | नियम फिल्टर करण्यासाठी |
| `n` | नवीन नियम तयार करण्यासाठी |
| `e` | निवडलेला नियम संपादित (edit) करण्यासाठी |
| `d` | निवडलेला नियम हटवण्यासाठी (delete) |
| `Space` | टनेल चालू/बंद (start/stop) टॉगल करण्यासाठी |
| `l` | निवडलेल्या टनेलसाठी लॉग पाहण्यासाठी |
| `g` | ग्लोबल (global) लॉग पाहण्यासाठी |
| `t` | थीम बदलण्यासाठी (toggle theme) |
| `?` | मदत (help) मेनू टॉगल करण्यासाठी |
| `q` | बाहेर पडण्यासाठी (Quit) |
| `Esc` | मागे जाण्यासाठी / रद्द करण्यासाठी |

## कॉन्फिगरेशन (Configuration)

नियम स्वयंचलितपणे जतन केले जातात. डीफॉल्टनुसार, कॉन्फिगरेशन `~/.config/dwaar/rules.yaml` मध्ये सेव्ह केले जाते (फक्त वापरकर्त्यासाठी `0600` परवानग्यांसह). मोठ्या प्रमाणात नियम आयात (import) किंवा सुधारित करण्यासाठी तुम्ही ही फाइल स्वतः (manually) संपादित करू शकता.

उदाहरण:

```yaml
version: 1
rules:
  - id: rule_1
    name: db-direct
    local_port: 5432
    remote_host: 10.0.0.5
    remote_port: 5432
    protocol: tcp
    auto_start: true
  - id: rule_2
    name: redis-via-bastion
    local_port: 6379
    remote_host: 10.0.0.6
    remote_port: 6379
    protocol: tcp
    auto_start: false
    via:
      host: bastion.example.com
      port: 22
      user: deploy
      identity_file: ~/.ssh/id_ed25519
      # password: hunter2  # identity_file ऐवजी पर्याय
```

> **टीप (Note)**: कॉन्फिग फाइलमध्ये सेव्ह केलेला SSH पासवर्ड प्लेनटेक्स्टमध्ये (plaintext) ठेवला जातो, फक्त फाइल परवानग्यांनी संरक्षित. शक्य असल्यास की-आधारित (key-based) प्रमाणीकरण वापरा.

## विकास (Development)

प्रोजेक्टची रचना मानक Go प्रोजेक्ट लेआउट्सनुसार केली आहे:
- `cmd/dwaar/main.go` - ॲप्लिकेशनचा एंट्रीपॉइंट.
- `internal/tunnel` - टनेल व्यवस्थापन आणि SSH फॉरवर्डिंग लॉजिक.
- `internal/config` - कॉन्फिगरेशन हाताळणी.
- `internal/tui` - Bubble Tea व्ह्यूज, कंपोनंट्स आणि स्टायलिंग.

टेस्ट सूट चालवण्यासाठी:

```bash
go test ./...
```

## वेबसाइट (Website)

एक द्विभाषिक (इंग्रजी / मराठी) मार्केटिंग साइट [`site/`](site/) मध्ये आहे. हे
कोणत्याही अवलंबनाशिवाय स्थिर पृष्ठ आहे — `site/index.html` थेट उघडा, किंवा
फोल्डर सर्व्ह करा:

```bash
cd site
python3 -m http.server 8080
# नंतर http://localhost:8080 उघडा
```

अधिक माहितीसाठी [`site/README.md`](site/README.md) पाहा.

## परवाना (License)

[MIT License](https://opensource.org/licenses/MIT)
