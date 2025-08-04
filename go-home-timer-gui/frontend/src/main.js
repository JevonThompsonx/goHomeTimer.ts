import './style.css';
import {CalculateDepartureTime, CreateCalendarEvent} from '../wailsjs/go/main/App';

document.querySelector('#app').innerHTML = `
  <div class="container">
    <h1>Go Home Timer</h1>
    <div class="input-group">
      <label for="arrival-time">Arrival Time (HH:MM AM/PM):</label>
      <input id="arrival-time" type="text" />
    </div>
    <div class="input-group">
      <label>Break Duration:</label>
      <div class="break-buttons">
        <button id="break-30">30 mins</button>
        <button id="break-60">1 hour</button>
      </div>
    </div>
    <button id="calculate-btn">Calculate & Create Event</button>
    <div id="result"></div>
  </div>
`;

const arrivalTimeInput = document.getElementById('arrival-time');
const break30Btn = document.getElementById('break-30');
const break60Btn = document.getElementById('break-60');
const calculateBtn = document.getElementById('calculate-btn');
const resultDiv = document.getElementById('result');

let breakDuration = '0.5'; // Default break duration

break30Btn.addEventListener('click', () => {
  breakDuration = '0.5';
  break30Btn.classList.add('selected');
  break60Btn.classList.remove('selected');
});

break60Btn.addEventListener('click', () => {
  breakDuration = '1';
  break60Btn.classList.add('selected');
  break30Btn.classList.remove('selected');
});

calculateBtn.addEventListener('click', () => {
  const arrivalTime = arrivalTimeInput.value.trim();
  if (!arrivalTime) {
    resultDiv.innerText = 'Please enter an arrival time.';
    return;
  }

  CalculateDepartureTime(arrivalTime, breakDuration)
    .then(departureTime => {
      resultDiv.innerText = `Recommended Departure Time: ${departureTime}`;
      return CreateCalendarEvent(departureTime);
    })
    .then(eventResult => {
      resultDiv.innerText += `\n${eventResult}`;
    })
    .catch(err => {
      resultDiv.innerText = `Error: ${err}`;
    });
});