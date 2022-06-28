import {useSpring, animated} from 'react-spring'
import React from "react";

function Introduction() {
    const props = useSpring({
        to: {opacity: 1},
        from: {opacity: 0},
        reset: true,
        delay: 400,
    })
    return (
        <div className="absolute left-0 right-0 top-32">
            <animated.h1 className="text-center" style={props}>
                Valitse mieltymyksiesi mukaan
            </animated.h1>
        </div>
    )
}

export default Introduction;