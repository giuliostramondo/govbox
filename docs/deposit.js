var chartDom = document.getElementById('main');
var myChart = echarts.init(chartDom);
var option;

let nbBlocks = 0
let defaultDeposit=250
let deposits = [];
let numProposals = []
let N = 2
for (var i = 0; i < nbBlocks; i++) {
  deposits.push({value:[i,defaultDeposit]});
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
			max: 6,
    },
	],
	legend: {
    data: ['Num proposals', 'Deposit'],
  }
};

function reset() {
	deposits = [];
	numProposals = [];
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
	let lastDeposit = defaultDeposit
	if (deposits.length>0) {
		lastDeposit = deposits[deposits.length-1].value[1]
	}
	let alpha = parseFloat(document.getElementById('alphaUp').value);
	let beta = 0
	if (n <= N) {
		alpha = parseFloat(document.getElementById('alphaDown').value);
		beta = -1
	}
	let k = parseFloat(document.getElementById('k').value);
  let v = (Math.abs(n - N + beta)) ** (1/k)

	let D = lastDeposit * (1 + alpha * v)
	console.log(`lastDeposit=${lastDeposit} alpha=${alpha} beta=${beta} k=${k} n=${n}`)
	console.log(`n - N + beta = ${n - N + beta}`)
	console.log(`(Math.abs(n - N + beta)) ** (1/k) = ${v}`)
	console.log(`alpha * (Math.abs(n - N + beta)) ** (1/k) = ${alpha * v}`)
	if (D < defaultDeposit) {
		D = defaultDeposit
	}
	return D
}

setInterval(function () {
	if (paused) {
		return
	}
	nbBlocks++
	n = parseInt(document.getElementById('numProposals').value);
  // deposits.shift();
  deposits.push({ value:[ nbBlocks, computeDeposit(n) ] });
	// numProposals.shift();
	numProposals.push({ value:[ nbBlocks, n ] })
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

