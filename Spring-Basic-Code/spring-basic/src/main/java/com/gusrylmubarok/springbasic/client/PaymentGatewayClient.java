package com.gusrylmubarok.springbasic.client;

import lombok.Data;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Data
public class PaymentGatewayClient {

    private String endpoint;

    private String privateKey;

    private String publicKey;
}
