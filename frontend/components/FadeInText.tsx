import {useSpring, animated} from 'react-spring'
import React from "react";

function FadeInText(text: String) {
    const props = useSpring({
        to: {opacity: 1},
        from: {opacity: 0},
        reset: true,
        delay: 400,
    })
    return (
        <div className="absolute left-0 right-0 top-32">
            <animated.h1 className="text-center text-xl" style={props}>
                {text}
            </animated.h1>
        </div>
    )
}

export default FadeInText;