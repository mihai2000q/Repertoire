import { render } from '@testing-library/react'
import CustomIconArrowLeft from './CustomIconArrowLeft.tsx'
import CustomIconArrowRight from './CustomIconArrowRight.tsx'
import CustomIconElectricGuitar from './CustomIconElectricGuitar.tsx'
import CustomIconElectricGuitarST from './CustomIconElectricGuitarST.tsx'
import CustomIconFlyingVGuitar from './CustomIconFlyingVGuitar.tsx'
import CustomIconFlyingV2Guitar from './CustomIconFlyingV2Guitar.tsx'
import CustomIconGuitarHead from './CustomIconGuitarHead.tsx'
import CustomIconLesPaulGuitar from './CustomIconLesPaulGuitar.tsx'
import CustomIconLesPaulGuitarOutlined from './CustomIconLesPaulGuitarOutlined.tsx'
import CustomIconMetronome from './CustomIconMetronome.tsx'
import CustomIconLightningTrio from './CustomIconLightningTrio.tsx'

describe.concurrent('Custom Icons', () => {
  it('should render Custom Icon Arrow Left', () => {
    render(<CustomIconArrowLeft />)
  })

  it('should render Custom Icon Arrow Right', () => {
    render(<CustomIconArrowRight />)
  })

  it('should render Custom Icon Electric Guitar', () => {
    render(<CustomIconElectricGuitar />)
  })

  it('should render Custom Icon Electric Guitar ST', () => {
    render(<CustomIconElectricGuitarST />)
  })

  it('should render Custom Icon Flying V2 Guitar', () => {
    render(<CustomIconFlyingV2Guitar />)
  })

  it('should render Custom Icon Flying V Guitar', () => {
    render(<CustomIconFlyingVGuitar />)
  })

  it('should render Custom Icon Guitar Head', () => {
    render(<CustomIconGuitarHead />)
  })

  it('should render Custom Icon Les Paul Guitar', () => {
    render(<CustomIconLesPaulGuitar />)
  })

  it('should render Custom Icon Les Paul Guitar Outlined', () => {
    render(<CustomIconLesPaulGuitarOutlined />)
  })

  it('should render Custom Icon Lightning Trio', () => {
    render(<CustomIconLightningTrio />)
  })

  it('should render Custom Icon Metronome', () => {
    render(<CustomIconMetronome />)
  })
})
