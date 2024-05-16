// daje random int uključujući min i max vrednosti
        function rndInt(min, max) {
          return Math.floor(Math.random() * (max - min + 1) ) + min;
        }

        // daje novi niz sa nasumičnim elementima iz celog niza34523523
        // menja niz koji je ubačen!!!
        const array_rand2 = (niz) => {
          let novi_niz = [0];
          let i = 0;
          let izabrani_clan;
          while (niz.length > 0) {
            izabrani_clan = rndInt(0, niz.length-1);
            novi_niz[i] = niz[izabrani_clan];
            niz.splice(izabrani_clan, 1);
            i++;
          }
          return novi_niz;
        }

        // daje nov niz sa nasumičnim elementima iz niza tako da ih ima koliko broj_clanova
        // PRAVI DUPLE I MORA DA SE SREDI PO UGLEDU NA splice MEHANIZAM ONOG GORE!!!!
        function array_rand(niz, broj_clanova) {
          //console.log('niz:' + niz, 'broj:' + broj_clanova);
          //console.log(niz.length);
          let novi_niz = [0];
          for (let i = 0; i < broj_clanova; i++) {
            let clan = rndInt(0, niz.length-1);
            //console.log('clan: ' + clan);
            //console.log(niz[clan]);
            novi_niz[i] = niz[clan];
            //console.log('novi niz for:'+novi_niz);
          }
          //console.log('novi niz out:'+novi_niz);
          return novi_niz;
        }

        //lažni shuffle samo za nizove sa 10 članova
        function shuffle() {
          let kombinacija = rndInt(0, 2);
          let novi_niz = [0];
          let kombinacije = [
            [4, 8, 7, 0, 2, 5, 3, 1, 9, 6],
            [6, 9, 1, 3, 5, 2, 0, 7, 8, 4],
            [2, 9, 3, 8, 6, 5, 1, 0, 4, 7]
          ];
          let izabrana = kombinacije[kombinacija];

          //for (let i = 0; i < 10; i++) {
            novi_niz = izabrana;
          //}
          return novi_niz;
        }



        function S1m_kombi (props) {

          //console.log('console kombi', props)

          let s1m_kombi = <>{props.novo_pitanje.kombi}</>

          return (
            s1m_kombi
          )
        }



        function Zadaci_s1 (props) {

          //console.log('PROPOVI: ', props)

          const [zadatak, setZadatak] = React.useState('s1m_kombi');

          const [novo_pitanje, setNovo_pitanje] = React.useState({
            a: 2,
            op2: '+',
            b: 2,
            t11: 'Марко',
            t22: 'Маша',
            t33: 'камиончића',
            kombi: 'Кликни на: Нови задатак!'
          });


          const promeniZadatak = (zadatak_pt) => {
            setZadatak(zadatak_pt);

            if (zadatak_pt == 'o1m_1_10') {
              setNovo_pitanje({
                a: 0,
                op2: '+',
                b: 0,
                t11: 'Милош',
                t22: 'Јелена',
                t33: 'колача',
                kombi: 'Кликни на: Нови задатак!'
              });
            } else {
              setNovo_pitanje({
                a: 0,
                op2: '*',
                b: 0,
                t11: 'Милош',
                t22: 'Јелена',
                t33: 'колача',
                kombi: 'Кликни на: Нови задатак!'
              });
            }



          }

          const novoPitanje = () => {

            //console.log(zadatak)

            switch (zadatak) {

              case 's1m_kombi':
                let kombi = rndInt(1, 3);
                //kombi = 2;
                //console.log('kombi:'+kombi);

                let c = '';
                let ponavljanje;
                let broj_elemenata;
                let cifara;

                //$svi = rand(0, 1) == 0 ? "prvih " . rand(cifara, 9) : "svih";
                //$raspored = rand(0, 1) == 0 ? "nije bitan" : "jeste bitan";
                //ponavljanje = rand(0, 1) == 0 ? "ne ponavljaju" : "ponavljaju";

                //$rez;
                let pzpn;
                let html = '';

                // echo "svi: " . $svi_elementi . ", raspored: ". bitan_raspored . ", ponavljanje: " . ponavljanje . "<br>";
                // print_r($e);

                switch (kombi) {

                  case 1: //permutacije
                    broj_elemenata = rndInt(5, 9);
                    cifara = broj_elemenata;

                    if (rndInt(0, 1) == 0) { //random brojevi sa ponavljanjem
                      let c1 = [ //cifre
                        rndInt(0, 9),
                        rndInt(0, 9),
                        rndInt(0, 9),
                        rndInt(0, 9),
                        rndInt(0, 9),
                        rndInt(0, 9),
                        rndInt(0, 9),
                        rndInt(0, 9),
                        rndInt(0, 9),
                        rndInt(0, 9)
                      ];
                      //console.log('c1:' + c1.length + ' ' + c1);
                      c = array_rand(c1, broj_elemenata);
                      let da_se_kopira = rndInt(0, c.length-1);
                      for (let i = 0; i < rndInt(1, 2); i++) { c[i] = c[da_se_kopira]; }
                      ponavljanje = "ponavljaju";
                      //console.log('c1: '+ c1);
                    } else { //random brojevi sa ponavljanjem
                      // let c2 = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
                      // shuffle(c2);
                      // js analogija za php shuffle ne radi pa mora lažni
                      let c2 = shuffle();
                      //console.log('c2: '+c2.length);
                      c = array_rand(c2, broj_elemenata);
                      ponavljanje = "ne ponavljaju";
                      //console.log('c2: '+ c2);
                    }

                    //console.log('c: '+ c);

                    let deljivih = "";
                    pzpn = "";
                    let tip = rndInt(1, 2);

                    //tip = 1; //////

                    switch (tip) {

                      case 1: //  prvi, zadnji, parni, neparni, deljivi

                        deljivih = rndInt(0, 1);
                        if (deljivih == 0) { // prva/zadnja n parna/neparna
                          deljivih = "";
                          pzpn = rndInt(0, 1) == 0 ? "ako su prva " + rndInt(2, 4) : "ako su zadnja " + rndInt(2, 4);
                          pzpn = rndInt(0, 1) == 0 ? pzpn + " broja parna i " : pzpn + " broja neparna i ";
                        }
                        else { // n cifreni deljivi sa n
                          deljivih = "deljivih sa " + rndInt(2, 5) + " ";
                        }

                        //html = "";

                        let td = [];

                        //console.log('c: ' + c.length);
                        for (let i = 0; i < c.length; i++) {
                          td.push(<td key={i} className="border border-black p-2">{c[i]}</td>)
                        }

                        html = <>
                        <br/>
                        Koliko {cifara} cifrenih brojeva {deljivih}
                        <br/>
                        može da se dobije od cifara
                        <br/>
                        <table><tbody><tr>{td}</tr></tbody></table>
                        {pzpn}ako se cifre {ponavljanje}?
                        <br/>
                        <br/>
                        <img src='/static/assignments/permutacije.gif' />
                        </>

                        /* + "<br>" +
                        "<audio controls autoplay>" +
                        "<source src='../Kenndog - Beethoven (Lyrics) if you see the homies with the guap.mp3' type='audio/mpeg'>" +
                        "</audio>";*/

                        break;

                      case 2: // tekstualni
                        let t1 = ["Na polici ", "U šupi", "U sobi", "U korpi", "U frižideru"];
                        let t2 = ["knjige", "mačke", "sveće", "jabuke", "veštice"];
                        let t11 = t1[rndInt(0, 4)];
                        let t22 = t2[rndInt(0, 4)];

                        let t3 = ["vatrene: ", "žute: ", "grozne: ", "plave: ", "glupe: ", "crvene: " , "divne: ", "zelene: "];
                        let t33 = "";
                        let b1 = rndInt(0, 7);
                        let b2 = rndInt(0, 7);
                        if (b1 == b2) {
                          if (b1 < 4) {
                          b2 = b1 + 1;
                          } else {
                          b2 = b1 - 1;
                          }
                        }
                        let t31 = t3[b1];
                        let t32 = t3[b2];
                        let b3 = rndInt(0, 7);
                        if (b3 == b1 || b3 == b2) {
                          for (let i = 0; i < 20; i = i + rndInt(1, 2)) {
                          let treci = rndInt(0, 5);
                          if (t3[treci] != t31 || t3[treci] != t32) {
                            t33 = t3[treci];
                          }
                          }
                        } else {
                          t33 = t3[b3];
                        }

                        html = '';

                        // html = <>
                        // <br/>Koliko {cifara} cifrenih brojeva {deljivih} može da se dobije od cifara
                        // <br/> <table><tbody><tr> {td} </tr></tbody></table>
                        // {pzpn} ako se cifre {ponavljanje} ?
                        // <br/><br/>
                        // <img src='../permutacije.gif' />
                        // </>

                        html = <><br/>
                        {t11} se nalaze {t22} sledećih boja:<br/>
                        {t31} {rndInt(1, 5)},&nbsp;
                        {t32} {rndInt(1, 5)}&nbsp;i&nbsp;
                        {t33} {rndInt(1, 5)}.<br/>
                        Na koliko načina se one mogu rasporediti tako da {t22} iste boje budu jedna do druge?
                        <br/><br/>
                        <img src='/static/assignments/patke.webp' /><br/>
                        </>  /* +
                        "<audio controls autoplay>" +
                        "<source src='../Rokeri s Moravu - Krkenzi kikiriki evri dej.mp3' type='audio/mpeg'>" +
                        "</audio>";*/
                        break;

                      default:
                        break;
                    }

                    break;

                  case 2:
                    //echo "varijacije";
                    broj_elemenata = rndInt(0, 9);
                    cifara = rndInt(3, 9);

                    if (rndInt(0, 1) == 0) { //rand brojevi sa ponavljanjem
                    let c1 = [ //cifre
                      rndInt(0, 9),
                      rndInt(0, 9),
                      rndInt(0, 9),
                      rndInt(0, 9),
                      rndInt(0, 9),
                      rndInt(0, 9),
                      rndInt(0, 9),
                      rndInt(0, 9),
                      rndInt(0, 9),
                      rndInt(0, 9)
                    ];
                    c = array_rand(c1, c1.length);
                    let da_se_kopira = rndInt(0, c.length-1); //radi dupliranja nekih random elemenata
                    for (let i = 0; i < rndInt(1, 2); i++) { c[i] = c[da_se_kopira]; }
                    ponavljanje = "ponavljaju";
                    } else {
                    // let c2 = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];
                    // shuffle(c2);
                    // js analogija za php shuffle ne radi pa mora lažni
                    let c2 = shuffle();
                    c = array_rand(c2, c2.length);
                    ponavljanje = "ne ponavljaju";
                    }

                    let pn = "";
                    pzpn = "";

                    if (rndInt(1, 2) == 1) { //samo parni/neparni ili prvih/zadnjih parnih/neparnih
                    pn = rndInt(0, 1) == 0 ? " parnih" : " neparnih";
                    } else {
                    pzpn = rndInt(0, 1) == 0 ? "ako su prva " + rndInt(2, 4) : "ako su zadnja " + rndInt(2, 4);
                    pzpn = rndInt(0, 1) == 0 ? pzpn + " broja parna i " : pzpn + " broja neparna i ";
                    }

                    let td2 = [];

                    //console.log('c: ' + c.length);
                    for (let i = 0; i < c.length; i++) {
                      td2.push(<td key={i} className="border border-black p-2">{c[i]}</td>)
                    }

                    html = "";
                    html = <>
                    <br/>
                    Koliko {cifara} cifrenih {pn} brojeva može da se dobije od cifara<br/>
                    <table><tbody><tr>{td2}</tr></tbody></table>
                    {pzpn} ako se cifre {ponavljanje}?
                    <br/><br/>
                    <img src='/static/assignments/varijacije.gif' /><br/></> /* +
                    "<audio controls autoplay>" +
                    "<source src='../Sammy K - Fatal Attraction (Lyrics) hell naw better believe i aint that one.mp3' type='audio/mpeg'>" +
                    "</audio>";*/

                    break;

                  case 3:
                    html = "";
                    html = <>
                      <br/>
                        Комбинације су у изради!
                      <br/>
                      <img src='/static/assignments/kombinacije.gif' />
                    </>
                    break;

                  default:
                    break;

                }

                //console.log('html: ', html)
                setNovo_pitanje({
                  a: 3,
                  op2: '+',
                  b: 3,
                  t11: 'Милош',
                  t22: 'Јелена',
                  t33: 'колача',
                  kombi: html
                });
                break;

              default:
                break;
            }

            // setVidljiv_odgovor(false);
            //setOdgovor(null);

          }


          return (
            <>
              <p className="text-2xl mt-5 ml-2">Одабери задатке:</p>

              <div className="mx-2 p-2 text-black max-w-sm border-2 border-sky-700 rounded-md bg-gradient-to-r from-blue-50 to-transparent">

                {/*
                <label className={zadatak == 'o1m_1_10'
                  ? "relative p-1 border-2 border-black rounded-md bg-sky-400"
                  : "relative p-1 border-2 border-black rounded-md bg-gray-400"}
                htmlFor="o1m_1_10">Сабирање и одузимање од 1 до 10
                </label>
                <input
                onClick={(e) => promeniZadatak('o1m_1_10')}
                className="relative ml-2 mt-2"
                type="radio" id="o1m_1_10" name="radio_z" value="o1m_1_10"
                defaultChecked
                />
                <br/>

                <label className={zadatak == 'o1m_1_10txt'
                  ? "relative p-1 border-2 border-black rounded-md bg-gradient-to-br from-red-500 via-yellow-100 to-green-300"
                  : "relative p-1 border-2 border-black rounded-md bg-gray-400"}
                htmlFor="o1m_1_10txt">Смешно саб. и одуз. од 1 до 10
                </label>
                <input
                onClick={(e) => promeniZadatak('o1m_1_10txt')}
                className="relative ml-2 mt-5"
                type="radio" id="o1m_1_10txt" name="radio_z" value="o1m_1_10txt"
                />
                <br/>
                */}

                <label className={zadatak == 's1m_kombi'
                  ? "relative p-1 border-2 border-black rounded-md bg-sky-400"
                  : "relative p-1 border-2 border-black rounded-md bg-gray-400"} htmlFor="s1m_kombi">Комбинаторика за средњу школу
                </label>
                <input
                onClick={(e) => promeniZadatak('s1m_kombi')}
                className="relative ml-2 mt-5"
                type="radio" id="s1m_kombi" name="radio_z" value="s1m_kombi"
                />
                <br/>

                {/*
                <label className={zadatak == 'o2m_mnozenje'
                  ? "relative p-1 border-2 border-black rounded-md bg-sky-400"
                  : "relative p-1 border-2 border-black rounded-md bg-gray-400"} htmlFor="o2m_mnozenje">Множење до 100
                </label>
                <input
                onClick={(e) => promeniZadatak('o2m_mnozenje')}
                className="relative ml-2 mt-5"
                type="radio" id="o2m_mnozenje" name="radio_z" value="o2m_mnozenje"
                />
                <br/>
                */}

                {zadatak == 's1m_kombi' && (
                <>
                  <br/>
                  <button onClick={novoPitanje}
                  className="ml-8 my-2 relative
                  border-2 border-gray-500 rounded-md p-1 bg-blue-300 hover:bg-blue-400">
                    Нови задатак</button>
                  <br/>

                  <S1m_kombi
                    novo_pitanje={novo_pitanje}
                  />
                </>
                )}

              </div>
              <br/>
            </>
          )
        }



        const container = document.getElementById('root');
        const root = ReactDOM.createRoot(container);
        root.render(<Zadaci_s1 />);
