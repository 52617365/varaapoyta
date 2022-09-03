import React, { useState, useEffect } from "react";
import { Countdown } from "react-daisyui";

function CountdownThingy({
  hours,
  minutes,
  seconds,
}: {
  hours: number;
  minutes: number;
  seconds: number;
}) {
  const [hoursState, setHours] = useState<number>(hours);
  const [minutesState, setMinutes] = useState<number>(minutes);
  const [secondsState, setSeconds] = useState<number>(seconds);

  useEffect(() => {
    if (hoursState === 0 && minutesState === 0 && secondsState === 0) {
      return;
    }

    if (secondsState === 0) {
      if (minutesState !== 0) {
        setMinutes(minutesState - 1);
        setSeconds(59);
      }
    }
    if (minutesState === 0) {
      if (hoursState !== 0) {
        setHours(hoursState - 1);
        setMinutes(59);
      }
    }
    if (hoursState === 0 && minutesState === 0 && secondsState === 0) {
      return;
    }
    const timer = setTimeout(() => {
      setSeconds((v: number) => (v <= 0 ? secondsState : v - 1));
    }, 1000);

    return () => {
      clearTimeout(timer);
    };
  }, [hoursState, minutesState, secondsState]);

  return (
    <div className="grid grid-flow-col gap-5 text-center auto-cols-max place-content-center">
      <div className="flex flex-col">
        <Countdown className="font-mono text-5xl" value={hoursState} />
        tunti
      </div>
      <div className="flex flex-col">
        <Countdown className="font-mono text-5xl" value={minutesState} />
        min
      </div>
      <div className="flex flex-col">
        <Countdown className="font-mono text-5xl" value={secondsState} />
        sek
      </div>
    </div>
  );
}

export default CountdownThingy;
