var chartDom = document.getElementById('main');
var myChart = echarts.init(chartDom);
var option;

let nbBlocks = 0
let baseDeposit = parseFloat(document.getElementById('baseDeposit').value);
let deposits = [];
let numProposals = []
for (var i = 0; i < nbBlocks; i++) {
  deposits.push({value:[i,baseDeposit]});
  numProposals.push({value:[i,0]});
}
option = {
  title: {
    text: 'Deposit throttling'
  },
  tooltip: {
    trigger: 'axis',
    formatter: function (params) {
			return 'height='+ params[0].data.value[0] +
				' n=' + params[1].data.value[1] +
				' D=' + params[0].data.value[1]
    },
    axisPointer: {
      animation: false
    }
  },
  xAxis: {
		name: 'Block height',
  },
  yAxis: [
		{
			name: 'Deposit',
    	type: 'value',
    	boundaryGap: [0, '25%'],
    },
		{
			name: 'Num proposals',
			min: 0,
			max: 8,
    },
	],
	legend: {
    data: ['Num proposals', 'Deposit'],
  }
};

function reset() {
	deposits = [];
	numProposals = [];
	nbBlocks=0;
	document.getElementById('numProposals').value = 0;
}

let paused = false

function pauseResume() {
	if (paused) {
		paused=false
		document.getElementById('pauseResume').innerHTML='Pause'
	} else {
		document.getElementById('pauseResume').innerHTML='Resume'
		paused=true
	}
}

function computeDeposit(n) {
	let baseDeposit = parseFloat(document.getElementById('baseDeposit').value);
	let lastDeposit = baseDeposit
	if (deposits.length>0) {
		lastDeposit = deposits[deposits.length-1].value[1]
	}
	let N = parseFloat(document.getElementById('N').value)
	let alpha = parseFloat(document.getElementById('alphaUp').value);
	if (n <= N) {
		alpha = parseFloat(document.getElementById('alphaDown').value);
	}
	let sign = 1
	if (n < N) {
		sign = -1
	}
	let k = parseFloat(document.getElementById('k').value);
  let v = (Math.abs(n - N)) ** (1/k)

	let D = lastDeposit * (1 + sign * alpha * v)
	console.log(`lastDeposit=${lastDeposit} alpha=${alpha} sign=${sign} k=${k} n=${n}`)
	console.log(`n - N = ${n - N}`)
	console.log(`(Math.abs(n - N)) ** (1/k) = ${v}`)
	console.log(`sign * alpha * (Math.abs(n - N)) ** (1/k) = ${sign * alpha * v}`)
	console.log(`----> D=${D}`)
	console.log('***********************************')
	if (D < baseDeposit) {
		D = baseDeposit
	}
	return D
}

function pushBlock(n) {
	nbBlocks++
  // deposits.shift();
  deposits.push({ value:[ nbBlocks, computeDeposit(n) ] });
	// numProposals.shift();
	numProposals.push({ value:[ nbBlocks, n ] })
}

setInterval(function () {
	if (paused) {
		return
	}
	blockWindow = parseInt(document.getElementById('blockWindow').value);
	n = parseInt(document.getElementById('numProposals').value);
	for (let i = 0; i < blockWindow; i++) {
		pushBlock(n)
	}
  myChart.setOption({
    series: [
			{
    	  type: 'line',
    	  data: deposits,
				symbolSize: 3,
    	},
    	{
    	  type: 'line',
    	  data: numProposals,
				yAxisIndex: 1,
				symbolSize: 3,
    	}
    ]
  });
}, 1000);

option && myChart.setOption(option);

