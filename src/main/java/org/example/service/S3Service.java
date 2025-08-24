package org.example.service;

import java.net.URL;
import java.time.Duration;

public interface S3Service {
    URL generateSignedUrl(String objectKey, Duration expiration);
}