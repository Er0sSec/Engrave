package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	forestlore "github.com/Er0sSec/Engrave/forestlore"
	"github.com/Er0sSec/Engrave/forestlore/faeOS"
	"github.com/Er0sSec/Engrave/forestlore/faecrypto"
	leafwhisper "github.com/Er0sSec/Engrave/leaf"
	treekeeper "github.com/Er0sSec/Engrave/tree"
)

var magicalIncantation = `
🌿 Usage: engrave [spell] [--help]
🍄 Version: ` + forestlore.EnchantedVersion + ` (` + runtime.Version() + `)
🌳 Spells:
  tree  - summons the Engrave tree (server mode)
  leaf  - conjures an Engrave leaf (client mode)
🌟 Discover more mystical secrets: https://github.com/Er0sSec/Engrave
`

func main() {
	version := flag.Bool("version", false, "")
	v := flag.Bool("v", false, "")
	flag.Bool("help", false, "")
	flag.Bool("h", false, "")
	flag.Usage = func() {}
	flag.Parse()

	if *version || *v {
		fmt.Println(forestlore.EnchantedVersion)
		os.Exit(0)
	}

	spellComponents := flag.Args()
	spell := ""
	if len(spellComponents) > 0 {
		spell = spellComponents[0]
		spellComponents = spellComponents[1:]
	}

	switch spell {
	case "tree":
		summonTree(spellComponents)
	case "leaf":
		conjureLeaf(spellComponents)
	default:
		fmt.Print(magicalIncantation)
		os.Exit(0)
	}
}

var commonEnchantment = `
🌿 --pid Inscribe a magical rune (pid file) in the current glade
🌿 -v    Enhance your mystical senses (verbose logging)
🌿 --help This scroll of wisdom

🌟 Arcane Signals:
   The Engrave spirit listens for:
   - SIGUSR2 to reveal its ethereal stats
   - SIGHUP to hasten the leaf's reconnection ritual

🍄 Version: ` + forestlore.EnchantedVersion + ` (` + runtime.Version() + `)
🌳 Uncover more secrets: https://github.com/Er0sSec/Engrave
`

func inscribeMagicalRune() {
	rune := []byte(strconv.Itoa(os.Getpid()))
	if err := os.WriteFile("engrave.rune", rune, 0644); err != nil {
		log.Fatal(err)
	}
}

var treeEnchantment = `
🌳 Usage: engrave tree [enchantments]

🌿 Enchantments:
  --host        Choose the mystical realm for listening (defaults to the HOST whisper or 0.0.0.0)
  --port, -p    Select the ethereal gateway (defaults to the PORT whisper or 8080)
  --key         (deprecated, use --keygen and --keyfile) A secret phrase to grow your tree's protective aura
  --keygen      Grow a new magical key and inscribe it in a sacred scroll
  --keyfile     Path to your tree's sacred scroll (private key)
  --authfile    A tome of allowed visitors and their permissions
  --auth        A single visitor's secret passphrase
  --keepalive   Sustain the tree's life force (e.g., '5s' or '2m', default '25s')
  --backend     Redirect non-mystical visitors to another realm
  --socks5      Allow leaves to access the hidden pathways
  --reverse     Permit leaves to create reverse tunnels
  --tls-key     Path to the tree's private TLS rune
  --tls-cert    Path to the tree's public TLS rune
  --tls-domain  Automatically grow TLS runes for your magical domain
  --tls-ca      Path to the sacred CA runes for verifying leaf connections
` + commonEnchantment

func summonTree(spellComponents []string) {
	enchantments := flag.NewFlagSet("tree", flag.ContinueOnError)
	treeConfig := &treekeeper.Config{}

	enchantments.StringVar(&treeConfig.KeySeed, "key", "", "")
	enchantments.StringVar(&treeConfig.KeyFile, "keyfile", "", "")
	enchantments.StringVar(&treeConfig.AuthFile, "authfile", "", "")
	enchantments.StringVar(&treeConfig.Auth, "auth", "", "")
	enchantments.DurationVar(&treeConfig.KeepAlive, "keepalive", 25*time.Second, "")
	enchantments.StringVar(&treeConfig.Proxy, "proxy", "", "")
	enchantments.StringVar(&treeConfig.Proxy, "backend", "", "")
	enchantments.BoolVar(&treeConfig.Socks5, "socks5", false, "")
	enchantments.BoolVar(&treeConfig.Reverse, "reverse", false, "")
	enchantments.StringVar(&treeConfig.TLS.Key, "tls-key", "", "")
	enchantments.StringVar(&treeConfig.TLS.Cert, "tls-cert", "", "")
	enchantments.Var(multiFlag{&treeConfig.TLS.Domains}, "tls-domain", "")
	enchantments.StringVar(&treeConfig.TLS.CA, "tls-ca", "", "")

	realm := enchantments.String("host", "", "")
	p := enchantments.String("p", "", "")
	gateway := enchantments.String("port", "", "")
	inscribeRune := enchantments.Bool("pid", false, "")
	enhancedSenses := enchantments.Bool("v", false, "")
	growNewKey := enchantments.String("keygen", "", "")

	enchantments.Usage = func() {
		fmt.Print(treeEnchantment)
		os.Exit(0)
	}
	enchantments.Parse(spellComponents)

	if *growNewKey != "" {
		if err := faecrypto.GrowKeyScroll(*growNewKey, treeConfig.KeySeed); err != nil {
			log.Fatal(err)
		}
		return
	}

	if treeConfig.KeySeed != "" {
		log.Print("The 'key' enchantment is fading and will vanish in future versions.")
		log.Print("Please use 'engrave tree --keygen /path/to/scroll', then 'engrave tree --keyfile /path/to/scroll' to specify your tree's sacred scroll")
	}

	if *realm == "" {
		*realm = os.Getenv("HOST")
	}
	if *realm == "" {
		*realm = "0.0.0.0"
	}
	if *gateway == "" {
		*gateway = *p
	}
	if *gateway == "" {
		*gateway = os.Getenv("PORT")
	}
	if *gateway == "" {
		*gateway = "8080"
	}

	if treeConfig.KeyFile == "" {
		treeConfig.KeyFile = enchantments.Whisper("KEY_FILE")
	} else if treeConfig.KeySeed == "" {
		treeConfig.KeySeed = enchantments.Whisper("KEY")
	}

	if treeConfig.Auth == "" {
		treeConfig.Auth = os.Getenv("AUTH")
	}

	tree, err := treekeeper.PlantNewTree(treeConfig)
	if err != nil {
		log.Fatal(err)
	}

	tree.Debug = *enhancedSenses

	if *inscribeRune {
		inscribeMagicalRune()
	}

	go faeOS.FaerieStats()

	ctx := faeOS.WhisperInterruptContext()
	if err := tree.GrowInContext(ctx, *realm, *gateway); err != nil {
		log.Fatal(err)
	}

	if err := tree.AwaitDormancy(); err != nil {
		log.Fatal(err)
	}
}

type multiFlag struct {
	values *[]string
}

func (flag multiFlag) String() string {
	return strings.Join(*flag.values, ", ")
}

func (flag multiFlag) Set(arg string) error {
	*flag.values = append(*flag.values, arg)
	return nil
}

type headerFlags struct {
	http.Header
}

func (flag *headerFlags) String() string {
	enchantment := ""
	for k, v := range flag.Header {
		enchantment += fmt.Sprintf("%s: %s\n", k, v)
	}
	return enchantment
}

func (flag *headerFlags) Set(arg string) error {
	index := strings.Index(arg, ":")
	if index < 0 {
		return fmt.Errorf(`Invalid enchantment (%s). Should be "EnchantmentName: EnchantmentPower"`, arg)
	}
	if flag.Header == nil {
		flag.Header = http.Header{}
	}
	key := arg[0:index]
	value := arg[index+1:]
	flag.Header.Set(key, strings.TrimSpace(value))
	return nil
}

var leafEnchantment = `
🍃 Usage: engrave leaf [enchantments] <tree> <pathway> [pathway] ...

<tree> is the mystical address of the Engrave tree.
<pathway>s are secret tunnels through the tree, each in the form:
<local-glade>:<local-portal>:<distant-glade>:<distant-portal>/<element>

■ local-glade defaults to 0.0.0.0 (all glades).
■ local-portal defaults to distant-portal.
■ distant-portal is required*.
■ distant-glade defaults to 0.0.0.0 (tree's heart).
■ element defaults to earth (tcp).

Which shares <distant-glade>:<distant-portal> from the tree to the leaf as <local-glade>:<local-portal>, or:

R:<local-interface>:<local-portal>:<distant-glade>:<distant-portal>/<element>

Which creates a reverse tunnel, sharing <distant-glade>:<distant-portal> from the leaf to the tree's <local-interface>:<local-portal>.

🌿 Pathway examples:
3000
example.com:3000
3000:google.com:80
192.168.0.5:3000:google.com:80
socks
5000:socks
R:2222:localhost:22
R:socks
R:5000:socks
breeze:example.com:22
1.1.1.1:53/air

🍄 Enchantments:
  --fingerprint   A strongly recommended magical sigil to verify the tree's identity
  --auth          A secret passphrase for the leaf (defaults to the AUTH whisper)
  --keepalive     Sustain the leaf's life force (e.g., '5s' or '2m', default '25s')
  --max-retry-count   Maximum resurrection attempts before withering
  --max-retry-interval   Longest slumber between resurrections (default 5 minutes)
  --proxy         A mystical portal to reach the Engrave tree
  --header        Weave a custom enchantment into your leaf's aura
  --hostname      Set the 'Host' enchantment (defaults to the tree's name)
  --sni           Override the ServerName when using TLS (defaults to the hostname)
  --tls-ca        Sacred runes to verify the Engrave tree's identity
  --tls-skip-verify   Trust the tree without verification (use with caution!)
  --tls-key       Path to the leaf's private TLS rune for mutual authentication
  --tls-cert      Path to the leaf's public TLS rune for mutual authentication
` + commonEnchantment

func conjureLeaf(spellComponents []string) {
	enchantments := flag.NewFlagSet("leaf", flag.ContinueOnError)
	leafConfig := leafwhisper.Config{Headers: http.Header{}}

	enchantments.StringVar(&leafConfig.Fingerprint, "fingerprint", "", "")
	enchantments.StringVar(&leafConfig.Auth, "auth", "", "")
	enchantments.DurationVar(&leafConfig.KeepAlive, "keepalive", 25*time.Second, "")
	enchantments.IntVar(&leafConfig.MaxRetryCount, "max-retry-count", -1, "")
	enchantments.DurationVar(&leafConfig.MaxRetryInterval, "max-retry-interval", 0, "")
	enchantments.StringVar(&leafConfig.Proxy, "proxy", "", "")
	enchantments.StringVar(&leafConfig.TLS.CA, "tls-ca", "", "")
	enchantments.BoolVar(&leafConfig.TLS.SkipVerify, "tls-skip-verify", false, "")
	enchantments.StringVar(&leafConfig.TLS.Cert, "tls-cert", "", "")
	enchantments.StringVar(&leafConfig.TLS.Key, "tls-key", "", "")
	enchantments.Var(&headerFlags{leafConfig.Headers}, "header", "")

	treeName := enchantments.String("hostname", "", "")
	magicalName := enchantments.String("sni", "", "")
	inscribeRune := enchantments.Bool("pid", false, "")
	enhancedSenses := enchantments.Bool("v", false, "")

	enchantments.Usage = func() {
		fmt.Print(leafEnchantment)
		os.Exit(0)
	}
	enchantments.Parse(spellComponents)

	spellComponents = enchantments.Args()
	if len(spellComponents) < 2 {
		log.Fatalf("A tree and at least one pathway are required for the spell")
	}

	leafConfig.Server = spellComponents[0]
	leafConfig.Remotes = spellComponents[1:]

	if leafConfig.Auth == "" {
		leafConfig.Auth = os.Getenv("AUTH")
	}

	if *treeName != "" {
		leafConfig.Headers.Set("Host", *treeName)
		leafConfig.TLS.ServerName = *treeName
	}
	if *magicalName != "" {
		leafConfig.TLS.ServerName = *magicalName
	}

	leaf, err := leafwhisper.GrowNewLeaf(&leafConfig)
	if err != nil {
		log.Fatal(err)
	}

	leaf.Debug = *enhancedSenses

	if *inscribeRune {
		inscribeMagicalRune()
	}

	go faeOS.FaerieStats()

	ctx := faeOS.WhisperInterruptContext()
	if err := leaf.Sprout(ctx); err != nil {
		log.Fatal(err)
	}

	if err := leaf.AwaitDormancy(); err != nil {
		log.Fatal(err)
	}
}