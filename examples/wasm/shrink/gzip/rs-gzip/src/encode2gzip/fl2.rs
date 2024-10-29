use std::io::Read;

use flate2::bufread::GzEncoder;

use super::enc2gz::CompressionLevel;

pub fn level2fl(l: CompressionLevel) -> flate2::Compression {
    match l {
        CompressionLevel::None => flate2::Compression::none(),
        CompressionLevel::Fast => flate2::Compression::fast(),
        CompressionLevel::Best => flate2::Compression::best(),
    }
}

pub struct Encoder {
    pub level: CompressionLevel,
    pub limit: u64,
}

impl Default for Encoder {
    fn default() -> Self {
        Self {
            level: CompressionLevel::default(),
            limit: super::enc2gz::GZ_SIZE_LIMIT,
        }
    }
}

impl Encoder {
    pub fn encode2gzip(&self, original: &[u8], out: &mut Vec<u8>) -> Result<(), &'static str> {
        super::enc2gz::encode2gzip(
            original,
            out,
            self.level,
            self.limit,
            &|org: &[u8], out: &mut Vec<u8>, l: CompressionLevel, lmt: u64| {
                let rdr = org;
                let lvl: flate2::Compression = level2fl(l);
                let gz = GzEncoder::new(rdr, lvl);
                let mut taken = gz.take(lmt);
                taken
                    .read_to_end(out)
                    .map_err(|_| "unable to read to end")?;
                Ok(())
            },
        )
    }
}
