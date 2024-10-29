pub const GZ_SIZE_LIMIT: u64 = 1_048_576;

#[derive(Clone, Copy)]
pub enum CompressionLevel {
    Fast,
    Best,
    None,
}

impl Default for CompressionLevel {
    fn default() -> Self {
        Self::Fast
    }
}

pub fn encode2gzip<E>(
    original: &[u8],
    out: &mut Vec<u8>,
    level: CompressionLevel,
    limit: u64,
    encoder: &E,
) -> Result<(), &'static str>
where
    E: Fn(&[u8], &mut Vec<u8>, CompressionLevel, u64) -> Result<(), &'static str>,
{
    out.clear();
    encoder(original, out, level, limit)
}
