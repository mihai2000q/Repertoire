import { forwardRef, SVGProps } from 'react'

interface CustomIconVoiceProps extends SVGProps<SVGSVGElement> {
  color?: string
  size?: number | string
}

const CustomIconVoice = forwardRef<SVGSVGElement, CustomIconVoiceProps>(
  ({ color, size = 24, ...props }, ref) => (
    <svg
      ref={ref}
      xmlns="http://www.w3.org/2000/svg"
      xmlnsXlink={'http://www.w3.org/1999/xlink'}
      width={size}
      height={size}
      fill={color || 'currentColor'}
      stroke={color || 'currentColor'}
      viewBox="0 0 64 64"
      preserveAspectRatio="xMidYMid meet"
      aria-hidden="true"
      {...props}
    >
      <path
        d="M38.478 42.632c-4.644-1.466-3.322-2.633 1.11-4.298c2.123-.799.832-2.484.89-3.832c.026-.617 2.452-.494 2.276-2.874c-.124-1.676-3.816-4.064-4.822-5.083c-.581-.588 1.184-2.197-.059-3.612c-1.697-1.934-1.965-5.299-2.992-7.181c0 0 .783-1.196.183-1.876c-5.176-5.859-24.491-5.321-29.427 3.302c-5.541 9.68-5.615 23.059 5.906 30.267C16.667 50.65 10.104 62 10.104 62h20.319c0-1.938-2.266-8.89 1.7-8.578c3.446.271 7.666.122 7.292-3.77c-.113-1.174-.246-2.231.574-3.204c.82-.972 2.007-2.706-1.511-3.816"
      ></path>
      <path d="M43.129 40.805L62 43.277v-4.943z"></path>
      <path d="M58.46 57.081l2.024-4.281l-17.355-9.368z"></path>
      <path d="M60.484 28.766l-2.024-4.282l-15.331 13.651z"></path>
    </svg>
  )
)

CustomIconVoice.displayName = 'CustomIconVoice'

export default CustomIconVoice
