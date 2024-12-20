interface CustomIconArrowRightProps {
  color?: string
  size?: number
}

function CustomIconArrowRight({ color, size = 24 }: CustomIconArrowRightProps) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      fill={color || "currentColor"}
      stroke={color || "currentColor"}
      strokeLinecap="round"
      strokeLinejoin="round"
      viewBox="0 0 1024 1024"
    >
      <path d="M419.3 264.8l-61.8 61.8L542.9 512 357.5 697.4l61.8 61.8L666.5 512z"/>
    </svg>
  )
}

export default CustomIconArrowRight
