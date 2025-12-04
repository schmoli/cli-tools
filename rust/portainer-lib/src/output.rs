use crate::error::PortainerError;
use serde::Serialize;

#[derive(Serialize)]
pub struct ErrorOutput {
    pub error: ErrorDetail,
}

#[derive(Serialize)]
pub struct ErrorDetail {
    pub code: String,
    pub message: String,
}

pub fn print_yaml<T: Serialize>(data: &T) -> Result<(), PortainerError> {
    let yaml = serde_yaml::to_string(data)
        .map_err(|e| PortainerError::ApiError(format!("Failed to serialize: {}", e)))?;
    print!("{}", yaml);
    Ok(())
}

pub fn print_error(err: &PortainerError) {
    let output = ErrorOutput {
        error: ErrorDetail {
            code: err.code().to_string(),
            message: err.to_string(),
        },
    };
    if let Ok(yaml) = serde_yaml::to_string(&output) {
        eprint!("{}", yaml);
    } else {
        eprintln!("error: {}", err);
    }
}
