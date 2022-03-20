package com.gusrylmubarok.springbasic.factory;

import com.gusrylmubarok.springbasic.client.PaymentGatewayClient;
import org.springframework.beans.factory.FactoryBean;
import org.springframework.stereotype.Component;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Component("paymentGatewayClient")
public class PaymentGatewayClientFactoryBean implements FactoryBean<PaymentGatewayClient> {

    @Override
    public PaymentGatewayClient getObject() throws Exception {
        PaymentGatewayClient client = new PaymentGatewayClient();
        client.setEndpoint("https://example.com");
        client.setPrivateKey("private");
        client.setPublicKey("public");
        return client;
    }

    @Override
    public Class<?> getObjectType() {
        return PaymentGatewayClient.class;
    }
}
