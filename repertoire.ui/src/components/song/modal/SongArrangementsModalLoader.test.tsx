import { mantineRender } from '../../../test-utils.tsx'
import SongArrangementsModalLoader from './SongArrangementsModalLoader.tsx'

describe('Song Arrangements Modal Loader', () => {
  it('should render', () => {
    mantineRender(<SongArrangementsModalLoader opened={true} onClose={vi.fn()} />)
  })
})
