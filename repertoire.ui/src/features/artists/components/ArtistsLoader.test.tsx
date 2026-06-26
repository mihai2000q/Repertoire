import { mantineRender } from '../../../test-utils.tsx'
import ArtistsLoader from './ArtistsLoader.tsx'

describe('Artists Loader', () => {
  it('should render', () => {
    mantineRender(<ArtistsLoader />)
  })
})
