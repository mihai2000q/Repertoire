import { mantineRender } from "../../../test-utils"
import ArtistSongsLoader from "./ArtistSongsLoader"

describe('Artist Songs Loader', () => {
  it('should render', () => {
    mantineRender(<ArtistSongsLoader />)
  })
})
