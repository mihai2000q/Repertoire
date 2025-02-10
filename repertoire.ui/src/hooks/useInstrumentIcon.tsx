import { ReactNode } from 'react'
import CustomIconVoice from '../components/@ui/icons/CustomIconVoice.tsx'
import CustomIconPiano from '../components/@ui/icons/CustomIconPiano.tsx'
import CustomIconAcousticGuitar from '../components/@ui/icons/CustomIconAcousticGuitar.tsx'
import CustomIconPianoKeyboard from '../components/@ui/icons/CustomIconPianoKeyboard.tsx'
import CustomIconDrums from '../components/@ui/icons/CustomIconDrums.tsx'
import CustomIconViolin from '../components/@ui/icons/CustomIconViolin.tsx'
import CustomIconSaxophone from '../components/@ui/icons/CustomIconSaxophone.tsx'
import CustomIconFlute from '../components/@ui/icons/CustomIconFlute.tsx'
import CustomIconHarp from '../components/@ui/icons/CustomIconHarp.tsx'
import CustomIconUkulele from '../components/@ui/icons/CustomIconUkulele.tsx'
import CustomIconBass from '../components/@ui/icons/CustomIconBass.tsx'
import CustomIconKingVGuitar from '../components/@ui/icons/CustomIconKingVGuitar.tsx'
import CustomIconTriangleMusic from '../components/@ui/icons/CustomIconTriangleMusic.tsx'
import { Instrument } from '../types/models/Song.ts'

const instrumentIcons = new Map<string, ReactNode>([
  ['Voice', <CustomIconVoice key={'voice'} size={'100%'} />],
  ['Piano', <CustomIconPiano key={'piano'} size={'100%'} />],
  ['Keyboard', <CustomIconPianoKeyboard key={'keyboard'} size={'100%'} />],
  ['Electric Guitar', <CustomIconKingVGuitar key={'electric-guitar'} size={'100%'} />],
  ['Acoustic Guitar', <CustomIconAcousticGuitar key={'acoustic-guitar'} size={'100%'} />],
  ['Bass', <CustomIconBass key={'bass'} size={'100%'} />],
  ['Ukulele', <CustomIconUkulele key={'ukulele'} size={'100%'} />],
  ['Drums', <CustomIconDrums key={'drums'} size={'100%'} />],
  ['Violin', <CustomIconViolin key={'violin'} size={'100%'} />],
  ['Saxophone', <CustomIconSaxophone key={'saxophone'} size={'100%'} />],
  ['Flute', <CustomIconFlute key={'flute'} size={'100%'} />],
  ['Harp', <CustomIconHarp key={'harp'} size={'100%'} />]
])

export default function useInstrumentIcon() {
  function getInstrumentIcon(instrument: string | null | undefined | Instrument): ReactNode {
    const instrumentName: string | null =
      (instrument as Instrument)?.name !== undefined
        ? (instrument as Instrument).name
        : typeof instrument === 'string'
          ? instrument
          : null
    return instrumentIcons.get(instrumentName) ?? <CustomIconTriangleMusic size={'100%'} />
  }

  return getInstrumentIcon
}
