// daje random int uključujući min i max vrednosti
function rndInt(min, max) {
	return Math.floor(Math.random() * (max - min + 1) ) + min;
}

// daje novi niz sa nasumičnim elementima iz celog niza
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

//lažni shuffle samo za nizove sa 10 članova  gsag
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


function O1m_1_10 (props) {

	let pitanje = 'neko pitanje';

	// console.log('props O1m_1_10', props)

	let a = props.novo_pitanje.a;
	let op2 = props.novo_pitanje.op2;
	let b = props.novo_pitanje.b;

	pitanje =
	<>
		<img className="mb-3" width='525' src='/static/assignments/brojevna-prava2.png'/>
		<label className='text-2xl ml-2' htmlFor="r">{a} {op2} {b} = </label>
		<input onChange={(e) => props.setOdgovor(e.target.value)} className='text-2xl w-16  bg-lime-300 px-1 border border-lime-700 rounded-md' type='number' min='0' max='10' id='r' name='r' defaultValue={0} />
		<span className='text-2xl ml-2 font-bold text-green-700 '>&nbsp;&nbsp;{props.resenje}</span>
	</>

	let o1m_1_10 = <>

		<div
		className="p-1 border border-black rounded-md bg-gradient-to-t from-lime-100 via-white to-white shadow-xl">
				{pitanje}
		</div>

		<div hidden={!props.vidljiv_odgovor}
		className="p-1 h-fit border border-t-0 border-black rounded-md bg-lime-50 shadow-xl">
				{props.rezultat}
		</div>

	</>

	return (
		o1m_1_10
	)
}


function O1m_1_10txt (props) {

	let pitanje = 'neko pitanje';

	//console.log('props', props)

	let a = props.novo_pitanje.a;
	//let op2 = props.novo_pitanje.op2;
	let b = props.novo_pitanje.b;
	let t11 = props.novo_pitanje.t11;
	let t22 = props.novo_pitanje.t22;
	let t33 = props.novo_pitanje.t33;

	let c;
	let label;


	if (a > b) {
			c = a - b;
			label = t11 + " има " + a + " " + t33 + " a " + t22 + " " + c + " мање. Колико имају заједно?";
		} else {
			c = b - a;
			label = t11 + " има " + a + " " + t33 + " a " + t22 + " " + c + " више. Колико имају укупно?";
		}

	pitanje =
	<>
		<img className="mb-3" width='525' src='/static/assignments/brojevna-prava2.png'/>
		<label className='text-2xl' htmlFor="r">{label} </label>
		<input onChange={(e) => props.setOdgovor(e.target.value)} className='text-2xl w-16 bg-lime-300 px-1 border border-lime-700 rounded-md' type='number' min='0' max='10' id='r' name='r' defaultValue={0} />
		<span className='text-2xl ml-2 font-bold text-green-700 '>&nbsp;&nbsp;{props.resenje}</span>
	</>

	let o1m_1_10txt = <>

		<div
		className="p-1 border border-black rounded-md bg-gradient-to-t from-lime-100 via-white to-white shadow-xl">
				{pitanje}
		</div>

		<div hidden={!props.vidljiv_odgovor}
		className="p-1 h-fit border border-t-0 border-black rounded-md bg-lime-50 shadow-xl">
				{props.rezultat}
		</div>

	</>

	return (
		o1m_1_10txt
	)
}

export function reactRenderZadaci_o1(tekst) {
	const rootElement = document.getElementById("root");
	if (!rootElement) {
		throw new Error(`Could not find element with iidd ${"root"}`);
	}
	const reactRoot = ReactDOM.createRoot(rootElement);
	reactRoot.render(<Zadaci_o1 tekst={tekst}/>);
}

export function Zadaci_o1 (props) {

	//console.log('PROPOVI: ', props)

	const [zadatak, setZadatak] = React.useState('o1m_1_10');
	const [vidljiv_odgovor, setVidljiv_odgovor] = React.useState(false);
	const [novo_pitanje, setNovo_pitanje] = React.useState({
		a: 2,
		op2: '+',
		b: 2,
		t11: 'Марко',
		t22: 'Маша',
		t33: 'камиончића',
		kombi: 'Кликни на: Нови задатак!'
	});

	const [odgovor, setOdgovor] = React.useState(null);
	const [resenje, setResenje] = React.useState(null);
	const [rezultat, setRezultat] = React.useState(null);

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


	const vidiOdgovor = () => {

		let rezultat;

		switch (zadatak) {
			case 'o1m_1_10':
				var {a, op2, b} = novo_pitanje;
				var r = odgovor;

				//console.log(a, op2, b, r)

				var rr = "";

				var tacno = 0;
				if (op2 == '-') {
					var rr = a - b;
				} else {
					var rr = a + b;
				}

				if (rr == r) {
				tacno = 1;
				} else {
				tacno = 0;
				}

				setResenje(rr);

				if (tacno == 1) {
					rezultat = <>
						<p style={{textAlign: "center", background: "skyblue", fontSize: "30px"}}>
							&#10004;
						</p>
						{/*
						<img src='/hepi.gif' />
						<audio controls autoPlay>
							<source src='/slavuj2.mp3' type='audio/mpeg' />
						</audio>
						*/}
					</>
				} else {
					rezultat = <>
						<p style={{background: "red", textAlign: "center", fontSize: "30px"}}>
							&#10008;
						</p>

						{/*
						<img src='/sad.gif' />
						<audio controls autoPlay>
							<source src='/beba.mp3' type='audio/mpeg' />
						</audio>
						*/}
					</>
				}
				break;

			case 'o1m_1_10txt':
				var {a, b} = novo_pitanje;
				var r = odgovor;

				var rr = a + b;
				var tacno = 0;

				if (rr == r) {
					tacno = 1;
				} else {
					tacno = 0;
				}

				setResenje(rr);

				if (tacno == 1) {
					rezultat = <>
						<p style={{textAlign: "center", background: "skyblue", fontSize: "30px"}}>
							&#10004;
						</p>
						{/*
						<img src='/hepi.gif' />
						<audio controls autoPlay>
							<source src='/slavuj2.mp3' type='audio/mpeg' />
						</audio>
						*/}
					</>
				} else {
					rezultat = <>
						<p style={{background: "red", textAlign: "center", fontSize: "30px"}}>
							&#10008;
						</p>

						{/*
						<img src='/sad.gif' />
						<audio controls autoPlay>
							<source src='/beba.mp3' type='audio/mpeg' />
						</audio>
						*/}
					</>
				}
				break;

			default:
		}



		if (vidljiv_odgovor == true) {
			setVidljiv_odgovor(false);
		} else {
			setVidljiv_odgovor(true);
			setRezultat(rezultat);
		}
		//console.log('odgovor: ', odgovor)
	}

	const novoPitanje = () => {

		//console.log(zadatak)

		switch (zadatak) {
			case 'o1m_1_10':
				let a = Math.floor(Math.random() * 11); //kao rand(0, 10); u php jer daje 0,324 islične brojeve
				let op1 = Math.floor(Math.random() * 2);
				let op2 = '';
				let b = Math.floor(Math.random() * 11);

				if (op1 == 0) {
					if (a < b) {
						let tmp = a;
						a = b;
						b = tmp;
					}
				} else {
					while (a+b > 10) {
						a = Math.floor(Math.random() * 11);
						b = Math.floor(Math.random() * 11);
					}
				}

				if (op1 == 0) {
					op2 = "-";
				} else {
					op2 = "+";
				}

				setResenje(null);

				setNovo_pitanje({
					a,
					op2,
					b,
					t11: 'Милош',
					t22: 'Јелена',
					t33: 'колача',
					kombi: 'Кликни на: Нови задатак!'
				});
				break;

			case 'o1m_1_10txt':
				let aa = rndInt(2, 10);
				let bb = rndInt(2, 10);

				let t1 = ["Марко", "Mајмун", "Милош", "Kрокодил", "Стефан", "Слепи миш", "Илија", "Слон", "Јован", "Горила"];

				let t2 = ["Маша", "жирафа", "Јелена", "хијена", "Ивана", "чапља", "Милена", "кокошка", "Наташа", "овца"];

				let t3 = ["камиончића", "ексера", "колача", "чекића", "динара", "фломастера" , "лизала", "цветића", "лептирића", "ћевапчића"];

				let t11 = t1[rndInt(0, 9)];
				let t22 = t2[rndInt(0, 9)];
				let t33 = t3[rndInt(0, 9)];

				while (aa+bb > 10) {
					aa = rndInt(2, 10);
					bb = rndInt(2, 10);
				}

				setResenje(null);

				setNovo_pitanje({
					a: aa,
					op2: '+', //ovde je nepotrebno ali da ne pravi neke bagove
					b: bb,
					t11,
					t22,
					t33,
					kombi: 'Кликни на: Нови задатак!'
				});
				break;

			default:
				break;
		}

		setVidljiv_odgovor(false);
		//setOdgovor(null);

	}




	return (
		<>

			{zadatak == 'o1m_1_10' && (
				<>
					
						<div>
							{props.id}
							{props.name}
							
							</div>
					
				</>
				)
			}

			<p id="naslov" className="text-2xl mt-5 ml-2">{props.tekst[0]}</p>

			<div className="mx-2 p-2 text-black max-w-sm border-2 border-sky-700 rounded-md bg-gradient-to-r from-blue-50 to-transparent">


				<label className={zadatak == 'o1m_1_10'
					? "relative p-1 border-2 border-black rounded-md bg-sky-400"
					: "relative p-1 border-2 border-black rounded-md bg-gray-400"}
				htmlFor="o1m_1_10">{props.tekst[1]}
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
				htmlFor="o1m_1_10txt">{props.tekst[2]}
				</label>
				<input
				onClick={(e) => promeniZadatak('o1m_1_10txt')}
				className="relative ml-2 mt-5"
				type="radio" id="o1m_1_10txt" name="radio_z" value="o1m_1_10txt"
				/>
				<br/>


				{zadatak == 'o1m_1_10' && (
				<>
					<br/>
					<button onClick={vidiOdgovor}
					className="ml-10 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-lime-50 hover:bg-lime-200" >
						{props.tekst[3]}</button>

					<button onClick={novoPitanje}
					className="ml-8 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-blue-300 hover:bg-blue-400">
						{props.tekst[4]}</button>
					<br/>

					<O1m_1_10
						novo_pitanje={novo_pitanje}
						vidljiv_odgovor={vidljiv_odgovor}
						setOdgovor={setOdgovor}
						rezultat={rezultat}
						resenje={resenje}
					/>
				</>
				)}

				{zadatak == 'o1m_1_10txt' && (
				<>

					<br/>
					<button onClick={vidiOdgovor}
					className="ml-10 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-lime-50 hover:bg-lime-200" >
						Решење</button>

					<button onClick={novoPitanje}
					className="ml-8 my-2 relative
					border-2 border-gray-500 rounded-md p-1 bg-blue-300 hover:bg-blue-400">
						Нови задатак</button>
					<br/>

					<O1m_1_10txt
						novo_pitanje={novo_pitanje}
						vidljiv_odgovor={vidljiv_odgovor}
						setOdgovor={setOdgovor}
						rezultat={rezultat}
						resenje={resenje}
					/>
				</>
				)}



			</div>
			<br/>
		</>
	)
}



// const container = document.getElementById('root');
// const root = ReactDOM.createRoot(container);
// root.render(<Zadaci_o1 />);
