var chartDom = document.getElementById('main');
var myChart = echarts.init(chartDom);
var option;

let nbBlocks = 2
let defaultFee=250
let fees = [];
let blockSize = []
let currentLearningRate = 0
let N = 2
for (var i = 0; i < nbBlocks; i++) {
  fees.push({value:[i,defaultFee]});
  blockSize.push({value:[i,0]});
}
option = {
  title: {
    text: 'Feemarket throttling'
  },
  tooltip: {
    trigger: 'axis',
    formatter: function (params) {
			return 'height='+ params[0].data.value[0] +
				' Block Size=' + params[1].data.value[1] +
				' Gas Price=' + params[0].data.value[1]
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
			name: 'Gas Price',
    	type: 'value',
    	boundaryGap: [0, '25%'],
    },
		{
			name: 'currentBlockSize',
			min: 0,
			max: 6,
    },
	],
	legend: {
    data: ['currentBlockSize', 'fee'],
  }
};

function reset() {
	fees = [];
	blockSize = []
	nbBlocks=0;
	document.getElementById('currentBlockSize').value = 3;
	document.getElementById('targetBlockSize').value = 3;
	document.getElementById('maxBlockSize').value = 6;
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

function computeFee(n) {
	let lastFee = defaultFee
	if (fees.length>0) {
		lastFee = fees[fees.length-1].value[1]
	}
	let alpha = parseFloat(document.getElementById('alpha').value);
	let beta = parseFloat(document.getElementById('beta').value);
	let delta = parseFloat(document.getElementById('delta').value);
	let window = parseFloat(document.getElementById('window').value);
	let gamma = parseFloat(document.getElementById('gamma').value);
	let currentBlockSize = parseFloat(document.getElementById('currentBlockSize').value);
	let targetBlockSize = parseFloat(document.getElementById('targetBlockSize').value);
	let maxBlockSize = parseFloat(document.getElementById('maxBlockSize').value);
	let maxLearningRate = parseFloat(document.getElementById('maxLearningRate').value);
	let minLearningRate = parseFloat(document.getElementById('minLearningRate').value);
	

	let beginningOfWindow = Math.max(0,blockSize.length-1 - window);
	console.log(`beginningOfWindow= ${beginningOfWindow}`)

	let blockConsumption = 0;
	for (var i = beginningOfWindow; 
		i < Math.min(beginningOfWindow + window, blockSize.length-1); 
		i++) 
	{
		blockConsumption += blockSize[i].value[1];
		console.log(`blockSize[${i}]= ${blockSize[i].value[1]}`)
	}
	console.log(`window = ${window}`)
	console.log(`maxBlockSize = ${maxBlockSize}`)
	blockConsumption /= window;
	blockConsumption /= maxBlockSize;
	console.log(`blockConsumption= ${blockConsumption}`)

	if(currentLearningRate == 0){
		currentLearningRate=minLearningRate;
	}
	var newLearningRate = Math.max(minLearningRate, beta*currentLearningRate);
	if(blockConsumption < gamma ||
		blockConsumption > 1-gamma){
		console.log(`blockConsumption < gamma => True`)
		newLearningRate = Math.min(maxLearningRate, alpha + currentLearningRate);
	}
	console.log(`newLearningRate = ${newLearningRate}`)
	currentLearningRate = newLearningRate;
	console.log(`currentBlockSize = ${currentBlockSize}`)
	console.log(`targetBlockSize = ${targetBlockSize}`)

	let netGasDelta = 0;
	let targetBlockUtilization = targetBlockSize/maxBlockSize;
	for (var i = beginningOfWindow;
		i < Math.min(beginningOfWindow + window, blockSize.length-1);
		i++)
	{
		let currentBlockUtilization = blockSize[i].value[1]/maxBlockSize;
		netGasDelta += currentBlockUtilization - targetBlockUtilization;
	}
	console.log(`netGasDelta = ${netGasDelta}`)
	let F = lastFee * (1 + (newLearningRate * (currentBlockSize - targetBlockSize))/targetBlockSize ) + delta * netGasDelta;
	console.log(`lastFee = ${lastFee}`)
	console.log(`currentFee = ${F}`)
	
	return F
}

setInterval(function () {
	if (paused) {
		return
	}
	nbBlocks++
	n = parseInt(document.getElementById('currentBlockSize').value);
  // deposits.shift();
  fees.push({ value:[ nbBlocks, computeFee(n) ] });
	// numProposals.shift();
	blockSize.push({ value:[ nbBlocks, n ] });
  myChart.setOption({
    series: [
			{
    	  type: 'line',
    	  data: fees,
				symbolSize: 3,
    	},
    	{
    	  type: 'line',
    	  data: blockSize,
				yAxisIndex: 1,
				symbolSize: 3,
    	}
    ]
  });
}, 1000);

option && myChart.setOption(option);

