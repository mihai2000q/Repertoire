import { render } from '@testing-library/react'
import CustomIconArrowLeft from './CustomIconArrowLeft.tsx'
import CustomIconArrowRight from './CustomIconArrowRight.tsx'
import CustomIconElectricGuitarWarlock from './CustomIconElectricGuitarWarlock.tsx'
import CustomIconKingVGuitar from './CustomIconKingVGuitar.tsx'
import CustomIconRhoadsGuitar from './CustomIconRhoadsGuitar.tsx'
import CustomIconGuitarHead from './CustomIconGuitarHead.tsx'
import CustomIconLesPaulGuitar from './CustomIconLesPaulGuitar.tsx'
import CustomIconLesPaulGuitarOutlined from './CustomIconLesPaulGuitarOutlined.tsx'
import CustomIconMetronome from './CustomIconMetronome.tsx'
import CustomIconLightningTrio from './CustomIconLightningTrio.tsx'
import CustomIconAcousticGuitar from './CustomIconAcousticGuitar.tsx'
import CustomIconFlute from './CustomIconFlute.tsx'
import CustomIconHarp from './CustomIconHarp.tsx'
import CustomIconPiano from './CustomIconPiano.tsx'
import CustomIconPianoKeyboard from './CustomIconPianoKeyboard.tsx'
import CustomIconSaxophone from './CustomIconSaxophone.tsx'
import CustomIconUkulele from './CustomIconUkulele.tsx'
import CustomIconViolin from './CustomIconViolin.tsx'
import CustomIconVoice from './CustomIconVoice.tsx'

describe.concurrent('Custom Icons', () => {
  it('should render Custom Icon Arrow Left', () => {
    render(<CustomIconArrowLeft />)
  })

  it('should render Custom Icon Arrow Right', () => {
    render(<CustomIconArrowRight />)
  })

  it('should render Custom Icon Acoustic Guitar', () => {
    render(<CustomIconAcousticGuitar />)
  })

  it('should render Custom Icon Electric Guitar Warlock', () => {
    render(<CustomIconElectricGuitarWarlock />)
  })

  it('should render Custom Icon Flute', () => {
    render(<CustomIconFlute />)
  })

  it('should render Custom Icon Guitar Head', () => {
    render(<CustomIconGuitarHead />)
  })

  it('should render Custom Icon Harp', () => {
    render(<CustomIconHarp />)
  })

  it('should render Custom Icon King V Guitar', () => {
    render(<CustomIconKingVGuitar />)
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

  it('should render Custom Icon Piano', () => {
    render(<CustomIconPiano />)
  })

  it('should render Custom Icon Piano Keyboard', () => {
    render(<CustomIconPianoKeyboard />)
  })

  it('should render Custom Icon Rhoads Guitar', () => {
    render(<CustomIconRhoadsGuitar />)
  })

  it('should render Custom Icon Saxophone', () => {
    render(<CustomIconSaxophone />)
  })

  it('should render Custom Icon Ukulele', () => {
    render(<CustomIconUkulele />)
  })

  it('should render Custom Icon Violin', () => {
    render(<CustomIconViolin />)
  })

  it('should render Custom Icon Voice', () => {
    render(<CustomIconVoice />)
  })
})
